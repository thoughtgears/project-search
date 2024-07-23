from google.cloud import aiplatform, storage
import os
import json

# Set environment variables to control logging
os.environ["GRPC_VERBOSITY"] = "ERROR"
os.environ["GRPC_TRACE"] = ""

project_id = os.getenv("GCP_PROJECT_ID")
region = os.getenv("GCP_REGION")
index_endpoint = "projects/924966064158/locations/europe-west2/indexEndpoints/3553621580972032"
index_id = "the_search_stream_1721689406544"
aiplatform.init(project=project_id, location=region)

client = storage.Client()
blob = client.get_bucket("project-search-vertex-vector").blob("0000048eaf718650477fec78a2416a86.json").download_as_string()
data = json.loads(blob)

endpoint = aiplatform.MatchingEngineIndexEndpoint(index_endpoint_name=index_endpoint, project=project_id)
response = endpoint.find_neighbors(
    deployed_index_id=index_id,
    queries=[data["image_embeddings"]],
    num_neighbors=10
)

print(response)
