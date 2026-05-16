# Chapter 12: Spark MLlib

## What is MLlib?

**MLlib** is Spark's scalable machine learning library. It runs on top of DataFrames and provides:

- **Feature engineering** (transformers)
- **ML algorithms** (estimators)
- **Pipelines** (chain transformers + estimators)
- **Model evaluation** and **hyperparameter tuning**
- **Model persistence** (save/load)

---

## The Pipeline API

MLlib follows the scikit-learn `Pipeline` pattern:

```
Transformer → transforms data (no fitting)
Estimator  → fits on data → produces a Transformer (the model)
Pipeline   → chains multiple Transformers and Estimators
```

```python
from pyspark.ml import Pipeline
from pyspark.ml.feature import (
    StringIndexer, OneHotEncoder, VectorAssembler,
    StandardScaler, Imputer
)
from pyspark.ml.classification import RandomForestClassifier
from pyspark.ml.evaluation import BinaryClassificationEvaluator

# 1. Define stages
indexer = StringIndexer(inputCol="category", outputCol="category_idx")
encoder = OneHotEncoder(inputCol="category_idx", outputCol="category_vec")
imputer = Imputer(inputCols=["age", "income"], outputCols=["age_imp", "income_imp"])

assembler = VectorAssembler(
    inputCols=["age_imp", "income_imp", "category_vec"],
    outputCol="features"
)

scaler = StandardScaler(inputCol="features", outputCol="scaled_features")
rf = RandomForestClassifier(featuresCol="scaled_features", labelCol="label", numTrees=100)

# 2. Build pipeline
pipeline = Pipeline(stages=[indexer, encoder, imputer, assembler, scaler, rf])

# 3. Train / Test split
train, test = df.randomSplit([0.8, 0.2], seed=42)

# 4. Fit
model = pipeline.fit(train)

# 5. Predict
predictions = model.transform(test)
predictions.select("label", "prediction", "probability").show()

# 6. Evaluate
evaluator = BinaryClassificationEvaluator(labelCol="label")
auc = evaluator.evaluate(predictions)
print(f"AUC: {auc:.4f}")

# 7. Save and load
model.save("models/rf_pipeline/")
from pyspark.ml import PipelineModel
loaded = PipelineModel.load("models/rf_pipeline/")
```

---

## Feature Engineering

### StringIndexer and OneHotEncoder

```python
# StringIndexer: "cat" → 0.0, "dog" → 1.0
indexer = StringIndexer(inputCol="animal", outputCol="animal_idx", handleInvalid="keep")
df = indexer.fit(df).transform(df)

# OneHotEncoder: 0.0 → [1,0], 1.0 → [0,1]
encoder = OneHotEncoder(inputCol="animal_idx", outputCol="animal_vec")
df = encoder.fit(df).transform(df)
```

### VectorAssembler

```python
assembler = VectorAssembler(
    inputCols=["feature1", "feature2", "feature3", "one_hot_vec"],
    outputCol="features",
    handleInvalid="skip"
)
df = assembler.transform(df)
```

### Scalers

```python
from pyspark.ml.feature import StandardScaler, MinMaxScaler, RobustScaler

# StandardScaler: zero mean, unit std
scaler = StandardScaler(inputCol="features", outputCol="scaled", withMean=True, withStd=True)
df = scaler.fit(df).transform(df)

# MinMaxScaler: scale to [0, 1]
scaler = MinMaxScaler(inputCol="features", outputCol="scaled", min=0.0, max=1.0)
```

### Text Features

```python
from pyspark.ml.feature import Tokenizer, HashingTF, IDF, Word2Vec, CountVectorizer

tokenizer = Tokenizer(inputCol="text", outputCol="words")

# TF-IDF
hashing   = HashingTF(inputCol="words", outputCol="tf", numFeatures=10000)
idf       = IDF(inputCol="tf", outputCol="tfidf")

# Word2Vec embeddings
w2v = Word2Vec(inputCol="words", outputCol="embedding", vectorSize=100, minCount=2)

# CountVectorizer (vocabulary-based)
cv = CountVectorizer(inputCol="words", outputCol="features", vocabSize=1000)
```

### Other Transformers

```python
from pyspark.ml.feature import (
    Bucketizer, QuantileDiscretizer,
    PCA, ChiSqSelector,
    SQLTransformer, Normalizer
)

# Bucketizer: bin continuous values
bucketizer = Bucketizer(
    splits=[0, 18, 25, 35, 60, float("inf")],
    inputCol="age", outputCol="age_bucket"
)

# PCA: dimensionality reduction
pca = PCA(k=50, inputCol="features", outputCol="pca_features")

# Feature selection with Chi-Square
selector = ChiSqSelector(numTopFeatures=20, featuresCol="features",
                          labelCol="label", outputCol="selected_features")

# SQL-based transformation
sql_t = SQLTransformer(statement="SELECT *, age * income AS age_income FROM __THIS__")
```

---

## Classification Algorithms

```python
from pyspark.ml.classification import (
    LogisticRegression,
    RandomForestClassifier,
    GBTClassifier,
    LinearSVC,
    NaiveBayes,
    DecisionTreeClassifier,
    MultilayerPerceptronClassifier,
)

# Logistic Regression
lr = LogisticRegression(
    featuresCol="features", labelCol="label",
    maxIter=100, regParam=0.01, elasticNetParam=0.8
)

# Random Forest
rf = RandomForestClassifier(
    featuresCol="features", labelCol="label",
    numTrees=200, maxDepth=10, featureSubsetStrategy="sqrt"
)

# Gradient Boosted Trees
gbt = GBTClassifier(
    featuresCol="features", labelCol="label",
    maxIter=50, maxDepth=5, stepSize=0.1
)

# Neural network (layers: input → hidden → output)
mlp = MultilayerPerceptronClassifier(
    featuresCol="features", labelCol="label",
    layers=[100, 64, 32, 2],    # input dim, hidden, output classes
    maxIter=100
)
```

---

## Regression Algorithms

```python
from pyspark.ml.regression import (
    LinearRegression,
    RandomForestRegressor,
    GBTRegressor,
    DecisionTreeRegressor,
    GeneralizedLinearRegression,
)

lr = LinearRegression(
    featuresCol="features", labelCol="price",
    regParam=0.1, elasticNetParam=0.5
)

rf_reg = RandomForestRegressor(
    featuresCol="features", labelCol="price",
    numTrees=100, maxDepth=8
)
```

---

## Clustering

```python
from pyspark.ml.clustering import KMeans, BisectingKMeans, GaussianMixture

# KMeans
kmeans = KMeans(featuresCol="features", k=5, seed=42, maxIter=20)
model  = kmeans.fit(df)
df     = model.transform(df)    # adds "prediction" column (cluster id)

# Evaluate (WSSSE — within-cluster sum of squared errors)
print(model.summary.trainingCost)

# Find optimal k with elbow method
costs = []
for k in range(2, 11):
    model = KMeans(k=k, seed=42).fit(df)
    costs.append((k, model.summary.trainingCost))
```

---

## Model Evaluation

```python
from pyspark.ml.evaluation import (
    BinaryClassificationEvaluator,
    MulticlassClassificationEvaluator,
    RegressionEvaluator,
    ClusteringEvaluator,
)

# Binary classification
bin_eval = BinaryClassificationEvaluator(
    labelCol="label", rawPredictionCol="rawPrediction",
    metricName="areaUnderROC"  # or "areaUnderPR"
)
auc = bin_eval.evaluate(predictions)

# Multiclass
mc_eval = MulticlassClassificationEvaluator(
    labelCol="label", predictionCol="prediction",
    metricName="accuracy"  # accuracy|f1|weightedPrecision|weightedRecall
)

# Regression
reg_eval = RegressionEvaluator(
    labelCol="price", predictionCol="prediction",
    metricName="rmse"   # rmse|mse|mae|r2
)
```

---

## Hyperparameter Tuning

```python
from pyspark.ml.tuning import ParamGridBuilder, CrossValidator, TrainValidationSplit

# Build grid
param_grid = ParamGridBuilder() \
    .addGrid(rf.numTrees, [50, 100, 200]) \
    .addGrid(rf.maxDepth, [5, 10, 15]) \
    .build()

# Cross-validation (k-fold)
cv = CrossValidator(
    estimator  = pipeline,
    estimatorParamMaps = param_grid,
    evaluator  = BinaryClassificationEvaluator(),
    numFolds   = 3,
    parallelism = 4,   # run N param sets in parallel
)
cv_model = cv.fit(train)
best_model = cv_model.bestModel

# Train-validation split (faster than CV)
tvs = TrainValidationSplit(
    estimator  = pipeline,
    estimatorParamMaps = param_grid,
    evaluator  = BinaryClassificationEvaluator(),
    trainRatio = 0.8,
)
tvs_model = tvs.fit(train)
```

---

## Feature Importance

```python
# After fitting a Random Forest model
rf_model = model.stages[-1]   # last stage in pipeline
importances = rf_model.featureImportances
feature_names = assembler.getInputCols()

for name, importance in sorted(zip(feature_names, importances), key=lambda x: -x[1]):
    print(f"{name}: {importance:.4f}")
```
