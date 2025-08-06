import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
import math

#functions to be applied
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
        'Entropy_of_sequence': get_entropy_of_items_array(items_array),
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

#dataset
df = pd.read_json("./train_data.json")

df["No_of_unique_items"] = df["Items"].apply(lambda x: get_unique_element(x))
df["Avg_item_freq"]=df["Items"].apply(lambda x: get_avg_item_freq(x))
df["Frequency_skewness"]=df["Items"].apply(lambda x: get_item_skewness(x))
df["Max_consecutive_duplicates"]=df["Items"].apply(lambda x: get_max_consecutive_duplicates(x))
df["Entropy_of_sequence"]=df["Items"].apply(lambda x: get_entropy_of_items_array(x))
df["Ratio_of_unique_items"]=df["Items"].apply(lambda x: get_unique_element(x)/len(x))
df["Cache_size"]=[25 for i in range(2000)]
df["Array_length"]=[1000 for i in range(2000)]

print(df.iloc[:1,[0]])
print(df.iloc[:,[0,2,4]])

s1=list(df.iloc[:,2])
s2=list(df.iloc[:,4])
arr=[]
for i in range(2000):
    arr.append(min(s1[i],s2[i]))

#defining best algo
best_algo=[]
for i in range(2000):
    if s1[i] == s2[i]:
        best_algo.append("NA")
    elif s1[i]==arr[i]:
        best_algo.append("LRU")
    elif s2[i]==arr[i]:
        best_algo.append("LFU")

# print(s1,s2)

df["Best_algo"]=best_algo
print(df["Best_algo"].value_counts())

#encoding best algo
from sklearn.preprocessing import LabelEncoder

label_encoder = LabelEncoder()
df['Best_algo_encoded'] = label_encoder.fit_transform(df['Best_algo'])

df["Median_reaccess_time"]=df["Items"].apply(lambda x:median_reaccess_time(x))
df["Percent_items_reused"]=df["Items"].apply(lambda x:percent_items_reused(x))
df["Std_dev_item_freq"]=df["Items"].apply(lambda x:std_dev_item_freq(x))
df["Run_length_entropy"]=df["Items"].apply(lambda x:run_length_entropy(x))
df["Longest_repeat_run_length"]=df["Items"].apply(lambda x:longest_repeat_run_length(x))
df["Gini_index_of_freq"]=df["Items"].apply(lambda x:gini_index_of_freq(x))
df["Item_position_variance"]=df["Items"].apply(lambda x:item_position_variance(x))
df["Reused_distance_variance"]=df["Items"].apply(lambda x:get_reuse_distance_variance(x))
df["Recency_frequency_ratio"]=df["Items"].apply(lambda x:get_recency_frequency_ratio(x))
# df=df.drop(['Items', 'LRUHits', 'LRUMiss', 'LFUHits', 'LFUMiss', 'Best_algo'],axis=1)
df.to_csv("train_data.csv")
