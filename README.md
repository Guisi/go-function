# Deploying to GCP

## Consumer

gcloud functions deploy posts_consumer --entry-point SavePost --runtime go113 --trigger-topic=posts --set-env-vars GOOGLE_CLOUD_PROJECT=guisi-portfolio

## Publisher
gcloud functions deploy posts_publisher --entry-point Publish --runtime go113 --trigger-http --set-env-vars GOOGLE_CLOUD_PROJECT=guisi-portfolio