## Deploying to GCP with Terraform

````
# set GCP credentials file var
set/export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/gcp_credentials.json

# init terraform
cd terraform & terraform init

# apply terraform changes
terraform apply -auto-approve -var 'project=<gcp_project_id>'

# test publisher http endpoint
curl -X POST "https://us-central1-<gcp_project_id>.cloudfunctions.net/posts_publisher" -H "Authorization: bearer $(gcloud auth print-identity-token)" -H "Content-Type:application/json" --data '{ "id": "1", "message": "Mensagem teste", "creationDate": "2021-04-02T18:40:32.000Z" }'
or use Postman collection inside load_tests folder

# destroy terraform changes
terraform destroy -auto-approve -var 'project=<gcp_project_id>'
````

## Deploying to GCP manually
#### Consumer
````
gcloud functions deploy posts_consumer --entry-point SavePost --runtime go113 --trigger-topic=posts --set-env-vars GOOGLE_CLOUD_PROJECT=<gcp_project_id>
````
#### Publisher
````
gcloud functions deploy posts_publisher --entry-point Publish --runtime go113 --trigger-http --set-env-vars GOOGLE_CLOUD_PROJECT=<gcp_project_id>
````