# Chapter 19: Contributing to Pimcore

## Overview

Contributing to Pimcore open source is a great way to level up your skills, build reputation, and give back to the community. This chapter walks through the entire contribution cycle — from finding an issue to getting your PR merged.

---

## 1. Find "Good First Issues"

- **GitHub:** `https://github.com/pimcore/pimcore`
- Go to **Issues** → **Labels** → `good first issue` or `help wanted`
- Pick small, focused issues: typo fix, docs improvement, small bug, minor feature

---

## 2. Fork & Clone

```bash
# Fork via GitHub UI, then:
git clone https://github.com/<your-username>/pimcore.git
cd pimcore

# Add upstream for syncing
git remote add upstream https://github.com/pimcore/pimcore.git
```

---

## 3. Set Up Core Dev Environment

```bash
# Install dependencies
composer install

# Copy test config
cp .env.test.dist .env.test

# Create test DB
mysql -u root -p -e "CREATE DATABASE pimcore_test;"

# Run DB schema init
bin/console doctrine:migrations:migrate --env=test --no-interaction

# Run full test suite to confirm everything works
vendor/bin/phpunit
```

---

## 4. Create a Branch & Make Changes

```bash
# Always branch from the current development branch (e.g., 10.x or 12.x)
git checkout -b fix/typo-in-product-service

# Make your changes
# ...

# Verify code quality
vendor/bin/phpunit
bin/console lint:container
bin/console lint:twig

# Commit with a clear message
git commit -m "Fix typo in ProductService docblock"
```

---

## 5. Submit a Pull Request

```bash
git push origin fix/typo-in-product-service
```

On GitHub → your fork → GitHub suggests creating a PR.
- **Base branch:** `pimcore/pimcore:10.x` (or current dev branch)
- **Compare branch:** your `fix/...`

---

## 6. PR Description Template

```markdown
## Changes in this PR
- Fixed a typo in ProductService docblock.
- Added missing test case for product slug generation.

## Why this is needed
Improves code readability and prevents confusion for new developers.

## Related Issues
Closes #12345

## Checklist
- [x] Tests pass locally (`vendor/bin/phpunit`)
- [x] No BC (backward compatibility) breaks
- [x] Code follows Pimcore coding style (PSR-12)
- [x] Docs updated if needed
```

---

## 7. After PR Submission

1. **CI runs automatically** → fix any failing checks
2. **Maintainers review** → may request changes → update your branch and push again
3. **Approved → Merged** 🎉

---

## Best Practices

```bash
# Keep branch in sync with upstream
git fetch upstream
git rebase upstream/10.x

# Fix code style
composer cs-fix   # or vendor/bin/php-cs-fixer fix

# Check for BC breaks
bin/console debug:container
```

- Always branch from `10.x` (not `main`)
- Keep commits small and focused
- Use **PSR-12 coding style**
- Write tests for your changes
- Keep PRs focused — one change per PR
