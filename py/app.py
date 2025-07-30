from flask import Flask, request, jsonify
import joblib
import numpy as np
import pandas as pd 

app = Flask(__name__)

MODEL_PATH = 'model.pkl'
SCALER_PATH = 'scaler.pkl'

try:
    model = joblib.load(MODEL_PATH)
    scaler = joblib.load(SCALER_PATH)
    print(f"Model loaded successfully from {MODEL_PATH}")
    print(f"Scaler loaded successfully from {SCALER_PATH}")
except Exception as e:
    print(f"ERROR: Could not load model or scaler. Please ensure '{MODEL_PATH}' and '{SCALER_PATH}' exist in the same directory as app.py. Error: {e}")
    model = None
    scaler = None 

@app.route('/predict', methods=['POST'])
def predict():
    """
    Predicts the best algorithm based on incoming JSON data.
    Expected JSON format:
    {
        "LFUAvgReaccess": float,
        "Frequency_skewness": float,
        "Entropy_of_sequence": float,
        "Ratio_of_unique_items": float,
        "Reused_distance_variance": float,
        "Recency_frequency_ratio": float
    }
    """
    if model is None or scaler is None:
        return jsonify({"error": "Model or Scaler not loaded on server. Please check server logs for details."}), 500

    try:
        data = request.get_json(force=True) 

        # Modify to what you want, and calculate yourself
        expected_features = [
            'Items',
            'AvgReaccess',
            'Miss',
            'Hit'
        ]

        input_values = []
        for feature_name in expected_features:
            if feature_name not in data:
                return jsonify({"error": f"Missing required feature: '{feature_name}' in the request data."}), 400
            input_values.append(data[feature_name])

        features_array = np.array(input_values).reshape(1, -1)

        scaled_features = scaler.transform(features_array)

        prediction_encoded = model.predict(scaled_features)[0] 

        response = {
            'prediction_encoded': int(prediction_encoded)         }
        return jsonify(response)

    except ValueError as e:
        return jsonify({"error": f"Invalid data type for one or more features. Please ensure all values are numeric. Details: {e}"}), 400
    except Exception as e:
        print(f"An unexpected error occurred: {e}") # Log the error on the server side
        return jsonify({"error": f"An internal server error occurred: {e}"}), 500

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)

