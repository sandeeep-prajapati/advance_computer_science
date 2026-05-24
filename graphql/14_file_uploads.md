# Chapter 14: File Uploads

## The GraphQL Multipart Spec

GraphQL doesn't natively support file uploads — files are binary data while GraphQL uses JSON. The **GraphQL multipart request spec** (by Jayden Seric) is the community standard for handling uploads.

The request is a `multipart/form-data` POST with three parts:
1. `operations` — the GraphQL query/mutation JSON
2. `map` — maps file variable paths to the multipart parts
3. File parts — the actual file data

---

## Server Setup (Apollo Server + graphql-upload)

```bash
npm install graphql-upload
```

```javascript
import { ApolloServer } from '@apollo/server';
import { graphqlUploadExpress } from 'graphql-upload';
import express from 'express';

const app = express();

// Must be before Apollo middleware
app.use(graphqlUploadExpress({
  maxFileSize: 10_000_000,   // 10 MB
  maxFiles: 5,
}));

// GraphQL Upload provides the Upload scalar
import { GraphQLUpload } from 'graphql-upload';

const resolvers = {
  Upload: GraphQLUpload,  // register the scalar
};
```

---

## Schema with Upload

```graphql
scalar Upload

type Mutation {
  uploadAvatar(file: Upload!): UploadResult!
  uploadFiles(files: [Upload!]!): [UploadResult!]!

  createPost(input: CreatePostWithImageInput!): Post!
}

input CreatePostWithImageInput {
  title: String!
  body: String!
  coverImage: Upload
}

type UploadResult {
  filename: String!
  mimetype: String!
  url: String!
  size: Int!
}
```

---

## Upload Resolver

The `Upload` scalar resolves to a Promise containing:
- `filename` — original file name
- `mimetype` — MIME type
- `encoding` — file encoding
- `createReadStream()` — function to create a readable stream

```javascript
import { createWriteStream, mkdir } from 'fs';
import { pipeline } from 'stream/promises';
import path from 'path';
import crypto from 'crypto';

const resolvers = {
  Mutation: {
    uploadAvatar: async (parent, { file }, context) => {
      requireAuth(context);

      const { filename, mimetype, createReadStream } = await file;

      // Validate file type
      const allowedTypes = ['image/jpeg', 'image/png', 'image/webp'];
      if (!allowedTypes.includes(mimetype)) {
        throw new GraphQLError('Invalid file type. Only JPEG, PNG, WebP allowed.', {
          extensions: { code: 'BAD_USER_INPUT' },
        });
      }

      // Generate unique filename
      const ext = path.extname(filename);
      const uniqueName = `${crypto.randomUUID()}${ext}`;
      const uploadDir = './uploads/avatars';
      const filePath = path.join(uploadDir, uniqueName);

      // Save to disk
      const stream = createReadStream();
      await pipeline(stream, createWriteStream(filePath));

      // Update user's avatar in DB
      const url = `/uploads/avatars/${uniqueName}`;
      await context.db.users.update(context.currentUser.id, { avatarUrl: url });

      return { filename: uniqueName, mimetype, url, size: 0 };
    },

    uploadFiles: async (parent, { files }, context) => {
      const uploads = await Promise.all(files.map(async (filePromise) => {
        const { filename, mimetype, createReadStream } = await filePromise;
        // ... save each file
        return { filename, mimetype, url: `/uploads/${filename}`, size: 0 };
      }));
      return uploads;
    },
  },
};
```

---

## Uploading to S3 (AWS)

In production, save to object storage — not local disk:

```javascript
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';

const s3 = new S3Client({ region: process.env.AWS_REGION });

async function uploadToS3(stream, filename, mimetype) {
  const key = `uploads/${Date.now()}-${filename}`;

  await s3.send(new PutObjectCommand({
    Bucket: process.env.S3_BUCKET,
    Key: key,
    Body: stream,
    ContentType: mimetype,
    ACL: 'public-read',
  }));

  return `https://${process.env.S3_BUCKET}.s3.amazonaws.com/${key}`;
}

const resolvers = {
  Mutation: {
    uploadAvatar: async (parent, { file }, context) => {
      requireAuth(context);
      const { filename, mimetype, createReadStream } = await file;

      // Validate
      if (!['image/jpeg', 'image/png', 'image/webp'].includes(mimetype)) {
        throw new GraphQLError('Invalid file type');
      }

      const stream = createReadStream();
      const url = await uploadToS3(stream, filename, mimetype);

      await context.db.users.update(context.currentUser.id, { avatarUrl: url });

      return { filename, mimetype, url, size: 0 };
    },
  },
};
```

---

## Image Processing (Sharp)

```javascript
import sharp from 'sharp';
import { Readable } from 'stream';

async function processAndUploadImage(createReadStream, filename) {
  const chunks = [];
  for await (const chunk of createReadStream()) {
    chunks.push(chunk);
  }
  const buffer = Buffer.concat(chunks);

  // Resize and convert to WebP
  const processed = await sharp(buffer)
    .resize(800, 800, { fit: 'inside', withoutEnlargement: true })
    .webp({ quality: 85 })
    .toBuffer();

  const processedStream = Readable.from(processed);
  const newFilename = filename.replace(/\.[^.]+$/, '.webp');

  const url = await uploadToS3(processedStream, newFilename, 'image/webp');
  return url;
}
```

---

## Client-Side Upload

### With Apollo Client (React):

```javascript
import { useMutation, gql } from '@apollo/client';

const UPLOAD_AVATAR = gql`
  mutation UploadAvatar($file: Upload!) {
    uploadAvatar(file: $file) {
      url
      filename
    }
  }
`;

function AvatarUpload() {
  const [uploadAvatar, { loading, error }] = useMutation(UPLOAD_AVATAR);

  const handleChange = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    try {
      const { data } = await uploadAvatar({
        variables: { file },
      });
      console.log('Uploaded:', data.uploadAvatar.url);
    } catch (err) {
      console.error('Upload failed:', err);
    }
  };

  return (
    <div>
      <input type="file" accept="image/*" onChange={handleChange} disabled={loading} />
      {loading && <p>Uploading...</p>}
      {error && <p>Error: {error.message}</p>}
    </div>
  );
}
```

Apollo Client automatically formats the multipart request — you just pass the `File` object as a variable.

### Raw fetch (without Apollo):

```javascript
async function uploadFile(file, query, variables = {}) {
  const formData = new FormData();

  formData.append('operations', JSON.stringify({
    query,
    variables: { ...variables, file: null },
  }));

  formData.append('map', JSON.stringify({ '0': ['variables.file'] }));
  formData.append('0', file);

  const response = await fetch('/graphql', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: formData,
  });

  return response.json();
}

// Usage
uploadFile(
  fileInputElement.files[0],
  `mutation UploadAvatar($file: Upload!) { uploadAvatar(file: $file) { url } }`
);
```

---

## Alternative: Pre-signed URLs (Better for Large Files)

For large files, generate a pre-signed S3 URL and let the client upload directly:

```graphql
type Mutation {
  getUploadUrl(filename: String!, contentType: String!): UploadUrlPayload!
  confirmUpload(key: String!): Post!
}

type UploadUrlPayload {
  uploadUrl: String!   # pre-signed S3 PUT URL
  key: String!         # S3 key to confirm after upload
}
```

```javascript
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { PutObjectCommand } from '@aws-sdk/client-s3';

const resolvers = {
  Mutation: {
    getUploadUrl: async (parent, { filename, contentType }, context) => {
      requireAuth(context);
      const key = `uploads/${context.currentUser.id}/${Date.now()}-${filename}`;

      const url = await getSignedUrl(s3, new PutObjectCommand({
        Bucket: process.env.S3_BUCKET,
        Key: key,
        ContentType: contentType,
      }), { expiresIn: 300 }); // URL valid for 5 minutes

      return { uploadUrl: url, key };
    },
  },
};
```

Client: PUT file directly to S3, then call `confirmUpload(key)` to save the record.

---

## Security Considerations

- Always validate MIME type server-side (don't trust client)
- Scan uploaded files for malware (ClamAV, AWS Macie)
- Set maximum file size limits
- Never serve user-uploaded files from the same origin as your app (use S3/CDN)
- Set Content-Disposition and Content-Type headers when serving
- Use random filenames — never trust user-provided filenames
- Store files outside the web root if serving locally
