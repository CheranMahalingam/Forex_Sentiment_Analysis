from technical_indicators import generate_technical_indicators
from twitter_sentiment import generate_twitter_sentiment
from store import store_indicator_data
from analysis_dataframe import get_technical_analysis_data
from preprocessing import preprocess
from predict import generate_prediction
from connections import scan_connections, broadcast_inference, store_inference

WINDOW_SIZE = 96
PREDICTION_TIME_STEPS = 4


def handler(event, context):
    try:
        date = event['Records'][0]['dynamodb']['Keys']['Date']['S']
        timestamp = event['Records'][0]['dynamodb']['Keys']['Timestamp']['S']

        analysis = generate_technical_indicators(event['Records'][0]['dynamodb']['NewImage'])
        sentiment = generate_twitter_sentiment(event['Records'][0]['dynamodb']['Keys'])
        store_indicator_data(date, timestamp, analysis, sentiment)
        indicator_df = get_technical_analysis_data(date, timestamp)
        processed_df = preprocess(indicator_df)
        predictions = generate_prediction("EURUSD", processed_df, WINDOW_SIZE, PREDICTION_TIME_STEPS, analysis["EURUSD"])
        inference_subscribers = scan_connections("EURUSD")
        broadcast_inference("EURUSD", inference_subscribers, predictions, date + timestamp)
        store_inference("EURUSD", predictions, timestamp)
    except Exception as e:
        print("Exception:", e)
