import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
import joblib

#dataset
df = pd.read_csv("train_data.csv")

df=df.drop(['Unnamed: 0'],axis=1)
print(df.info())

#independent features and dependent features
X=df[['LFUAvgReaccess','Frequency_skewness','Entropy_of_sequence','Ratio_of_unique_items','Reused_distance_variance','Recency_frequency_ratio']]
y=df["Best_algo_encoded"]

#train test split
from sklearn.model_selection import train_test_split,GridSearchCV
X_train,X_test,y_train,y_test= train_test_split(X,y,test_size=0.20)

#scaling the split
from sklearn.preprocessing import StandardScaler
scaler = StandardScaler()
X_train=scaler.fit_transform(X_train)
X_test = scaler.transform(X_test)   #don't use fit_transform for test data

#apply decision tree model
from sklearn.tree import DecisionTreeClassifier
# Define the parameter grid to search
param_grid = {
    'max_depth': [5, 7, 10, None], # None means unlimited depth
    'min_samples_split': [2, 5, 10],
    'min_samples_leaf': [1, 2, 4],
    'criterion': ['gini', 'entropy'],
    
}

# Initialize DecisionTreeClassifier
dtree = DecisionTreeClassifier()

# Initialize GridSearchCV
grid_search = GridSearchCV(estimator=dtree, param_grid=param_grid, cv=5, scoring='accuracy', n_jobs=-1, verbose=1)

# Fit GridSearchCV to the training data
grid_search.fit(X_train, y_train)
# Print the best parameters and best score
print(f"Best parameters: {grid_search.best_params_}")
print(f"Best cross-validation score: {grid_search.best_score_:.4f}")

# Evaluate the best model on the test set
best_model = grid_search.best_estimator_

#predict
y_pred=best_model.predict(X_test)

#check accuracy score
from sklearn.metrics import accuracy_score,classification_report
score=accuracy_score(y_pred,y_test)
print(score)
#classification report
print(classification_report(y_pred,y_test, zero_division=0))

model_filename = "model.pkl"
joblib.dump(best_model, model_filename)

scaler_filename = "scaler.pkl"
joblib.dump(scaler, scaler_filename)

