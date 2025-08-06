from flask import Flask, request, jsonify
import joblib
import numpy as np
import pandas as pd 


def gen_dict(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in gen_dict")
    d={}
    for i in items_array:
        if i not in d.keys():
            d[i]=1
        else:
            d[i]+=1
    return d

def get_unique_element(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in get_unique_element")
    if not items_array:
        return 0
    return len(list(gen_dict(items_array).keys()))

def get_avg_item_freq(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in get_avg_item_freq")
    if not items_array:
        return 0
    return np.mean(list(gen_dict(items_array).values()))

from scipy.stats import skew
def get_item_skewness(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in get_item_skewness")
    if not items_array:
        return 0
    frequencies = list(gen_dict(items_array).values())
    if np.std(frequencies) == 0:
        return 0 
    if len(frequencies)<3:
        return 0
    return skew(frequencies)

def get_max_consecutive_duplicates(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in max_consec_dupes")
    if not items_array:
        return 0
    max_consecutive = 0
    current_consecutive = 0
    for i in range(len(list(items_array))):
        if i == 0 or list(items_array)[i] != list(items_array)[i-1]:
            current_consecutive = 1
        else:
            current_consecutive += 1
        max_consecutive = max(max_consecutive, current_consecutive)
    return max_consecutive
def get_entropy_of_items_array(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in get_entropy_of_sequence")
    if not items_array:
        return 0
    item_counts = gen_dict(items_array)
    total_items = len(items_array)
    entropy = 0
    for count in item_counts.values():
        probability = count / total_items
        entropy -= probability * math.log2(probability)
    return entropy

def median_reaccess_time(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in median_reaccess_time")
    last_seen = {}
    reaccess_times = []
    for idx, item in enumerate(items_array):
        if item in last_seen:
            reaccess_times.append(idx - last_seen[item])
        last_seen[item] = idx
    return np.median(reaccess_times) if reaccess_times else 0

from collections import Counter

def percent_items_reused(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in percent_items_reused")
    freq = Counter(items_array)
    total_unique = len(freq)
    reused = sum(1 for item, count in freq.items() if count > 1)
    return reused / total_unique if total_unique else 0

def std_dev_item_freq(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in std_dev_item_freq")
    freq = list(Counter(items_array).values())
    return np.std(freq)

import math

def run_length_entropy(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in run_length_entropy")
    if not items_array:
        return 0
    run_lengths = []
    current_run = 1
    for i in range(1, len(items_array)):
        if items_array[i] == items_array[i-1]:
            current_run += 1
        else:
            run_lengths.append(current_run)
            current_run = 1
    run_lengths.append(current_run)
    probs = [r / sum(run_lengths) for r in run_lengths]
    return -sum(p * math.log2(p) for p in probs)

def longest_repeat_run_length(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in longest_repeat_run_length")
    max_run = current_run = 1
    for i in range(1, len(items_array)):
        if items_array[i] == items_array[i-1]:
            current_run += 1
            max_run = max(max_run, current_run)
        else:
            current_run = 1
    return max_run

def gini_index_of_freq(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in gini index of frequency")
    values = np.array(list(Counter(items_array).values()))
    sorted_vals = np.sort(values)
    n = len(sorted_vals)
    if n == 0:
        return 0
    cumulative_sum = np.cumsum(sorted_vals)
    return (n + 1 - 2 * np.sum(cumulative_sum) / cumulative_sum[-1]) / n

from collections import defaultdict

def item_position_variance(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in item_position_variance")
    pos_map = defaultdict(list)
    for idx, item in enumerate(items_array):
        pos_map[item].append(idx)
    variances = [np.var(positions) for positions in pos_map.values() if len(positions) > 1]
    return np.mean(variances) if variances else 0

def get_reuse_distance_variance(items_array):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in get_reuse_distance_variance")
    reuse_distances = []
    last_seen = {}
    for i, item in enumerate(items_array):
        if item in last_seen:
            reuse_distances.append(i - last_seen[item])
        last_seen[item] = i
    return np.var(reuse_distances) if reuse_distances else 0

def get_recency_frequency_ratio(items_array, window_size=50):
    if not isinstance(items_array, list):
        raise ValueError("Input must be a list in get_recency_frequency_ratio")
    recent_items = set(items_array[-window_size:])
    freq_dict = gen_dict(items_array)
    top_freq_items = set([k for k, v in sorted(freq_dict.items(), 
                         key=lambda x: x[1], reverse=True)[:window_size//2]])
    overlap = len(recent_items & top_freq_items)
    return overlap / min(len(recent_items), len(top_freq_items))
#main function
def extract_features(items_array):
    features_dict = {
        'No_of_unique_items': get_unique_element(items_array),
        'Avg_item_freq': get_avg_item_freq(items_array),
        'Frequency_skewness': get_item_skewness(items_array),
        'Max_consecutive_duplicates': get_max_consecutive_duplicates(items_array),
        'Entropy_of_items_array': get_entropy_of_items_array(items_array),
        'Ratio_of_unique_items': get_unique_element(items_array) / len(items_array) if len(items_array) > 0 else 0,
        'Median_reaccess_time': median_reaccess_time(items_array),
        'Percent_items_reused': percent_items_reused(items_array),
        'Std_dev_item_freq': std_dev_item_freq(items_array),
        'Run_length_entropy': run_length_entropy(items_array),
        'Longest_repeat_run_length': longest_repeat_run_length(items_array),
        'Gini_index_of_freq': gini_index_of_freq(items_array),
        'Item_position_variance': item_position_variance(items_array),
        'Reused_distance_variance': get_reuse_distance_variance(items_array),
        'Recency_frequency_ratio': get_recency_frequency_ratio(items_array),
    }
    return features_dict

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
        print('Items received:', data.get('Items'), type(data.get('Items')))
         
        # Modify to what you want, and calculate yourself
        expected_features = [
            'Items',
            'LFUAvgReaccess',
            'LFUHits',
            'LFUMiss',
            'LRUAvgReaccess',
            'LRUHits',
            'LRUMiss'  
        ]
        
        df = pd.read_csv("train_data.csv")
        for feature_name in expected_features:
            if (feature_name not in data) :
                data[feature_name] = df[feature_name].mean()   
            
        missing_keys = [key for key in expected_features if key not in data]
        if missing_keys:
            return jsonify({"error": f"Missing required keys in request data: {', '.join(missing_keys)}"}), 400
        # Type validations
        if not isinstance(data['Items'], list):
            return jsonify({"error": "Invalid type for 'Items'. Must be a list."}), 400
        for key in expected_features:
            if key != 'Items' and not isinstance(data[key], (int, float)):
                return jsonify({"error": f"Invalid type for '{key}'. Must be numeric."}), 400

        # features_array = np.array(input_values).reshape(1, -1)
        # print(features_array)
        features = extract_features(data["Items"])
        
        
        
        feature_vector = [
            data["LRUAvgReaccess"],
            data["LFUAvgReaccess"],
            features['Frequency_skewness'],
            features['Entropy_of_items_array'],  
            features['Ratio_of_unique_items'],
            features['Reused_distance_variance'],
            features['Recency_frequency_ratio']
        ]
          

        feature_vector= np.array(feature_vector).reshape(1,-1)
        scaled_vector = scaler.transform(feature_vector)
        
        prediction_encoded = model.predict(scaled_vector)[0] 

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

# 0 -> NA
# 1 -> LRU
# 2 -> LFU

