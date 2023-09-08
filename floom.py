import yaml
import argparse
import os
import requests
import glob
import json
import sys
import termcolor2
import colorama
from tqdm import tqdm
import time

colorama.init()

# Define a mapping to convert singular kinds to their plural API forms
kind_to_api_mapping = {
    'Pipeline': 'Pipelines',
    'Data': 'Data',
    'Embeddings': 'Embeddings',
    'Model': 'Models',
    'Prompt': 'Prompts',
    'Response': 'Responses',
    'VectorStore': 'VectorStores'
}

def load_config():
    with open("floom.yml", 'r') as stream:
        try:
            return yaml.safe_load(stream)
        except yaml.YAMLError as exc:
            print(exc)
            exit(1)

def get_engine(args, config):
    engine_name = args.engine or config['FloomEngines'][0]['name']
    for engine in config['FloomEngines']:
        if engine['name'] == engine_name:
            return engine
    print(f"Engine '{engine_name}' not found in configuration.")
    exit(1)

def upload_file(engine, filepath):
    url = f"{engine['url']}v1/Files"
    files = {'file': open(filepath, 'rb')}
    headers = {'Api-Key': f"{engine['apiKey']}"}
    
    print(termcolor2.colored(f"Uploading file '{filepath}'...", 'yellow'))
    
    response = requests.post(url, files=files, headers=headers)

    # Debugging lines
    #print(f"Status Code: {response.status_code}")
    #print(f"Response Text: {response.text}")

    if response.status_code == 200:
        #print (f"returning {json.loads(response.text)['fileId']}")
        return json.loads(response.text)['fileId']
    else:
        print("Failed to upload file.")
        exit(1)


def apply_yaml(engine, yaml_content):
    kind = yaml_content['kind']
    api_kind = kind_to_api_mapping.get(kind, kind)  # convert to plural/API format
    
    url = f"{engine['url']}v1/{api_kind}/Apply"
    headers = {
        'Content-Type': 'text/yaml',
        'Api-Key': f"{engine['apiKey']}"  
    }
    
    response = requests.post(url, data=yaml.dump(yaml_content), headers=headers)
    
    # Debugging lines
    # print(f"Status Code: {response.status_code}")
    # print(f"Response Text: {response.text}")
    
    return response.status_code == 200

def apply_file(args, config):
    file_path = args.file

    # Check if YAML file exists
    if not os.path.exists(file_path):
        print(termcolor2.colored(f"Error: YAML file '{file_path}' not found.", 'red'))
        exit(1)

    engine = get_engine(args, config)
    
   
    with open(file_path, 'r') as yaml_file:
        yaml_content = yaml.safe_load(yaml_file)

    if yaml_content['kind'] == 'Data' and yaml_content['type'] == 'file':
        if not os.path.exists(yaml_content['path']):
            print(f"File {yaml_content['path']} not found.")
            exit(1)
        file_id = upload_file(engine, yaml_content['path'])
        yaml_content['fileId'] = file_id
        yaml_content.pop('path', None)
        #print (yaml_content)

    if apply_yaml(engine, yaml_content):
        print(termcolor2.colored(f"{yaml_content['kind']} '{yaml_content['id']}' ('{file_path}') applied successfully.", 'green'))

    else:
        print("Failed to apply YAML.")


def apply_directory(args, config):
    directory = args.directory

    # Collect directory from arguments
    print(termcolor2.colored(f"Applying directory: {directory}",'yellow'))

    # Check if directory exists
    if not os.path.isdir(directory):
        print(termcolor2.colored(f"Error: Directory '{directory}' not found.", 'red'))
        exit(1)
    

    engine = get_engine(args, config)  # Assuming get_engine is a function you've defined elsewhere
    print(termcolor2.colored(f"Using engine '{engine['name']}'",'cyan'))

    # Pre-defined order of kinds
    order = ['VectorStore', 'Embeddings', 'Data', 'Model', 'Prompt', 'Response', 'Pipeline']



    # Initialize an empty list to store (kind, file) tuples
    files_to_apply = []

    # Scan all YAML files in the directory
    for yaml_file in glob.glob(os.path.join(directory, "*.yml")):
        with open(yaml_file, 'r') as file:
            content = yaml.safe_load(file)  # Load the YAML content
            kind = content.get('kind', None)  # Extract the 'kind' parameter
            obj_id = content.get('id', None)  # Extract the 'id' parameter

            if kind is not None and obj_id is not None:
                files_to_apply.append((kind, obj_id, yaml_file))

    # Sort the files according to the pre-defined 'order'
    files_to_apply.sort(key=lambda x: order.index(x[0]) if x[0] in order else len(order))

    # Apply files in sorted order
    for kind, obj_id, yaml_file in files_to_apply:
        print(termcolor2.colored(f"Applying {kind} '{obj_id}' ('{yaml_file}')...", 'yellow'))
        apply_file(argparse.Namespace(file=yaml_file, engine=engine['name']), config)
        
def list_pipelines(args, config):
    engine = get_engine(args, config)
    url = f"{engine['url']}v1/Pipelines/List"
    response = requests.get(url)
    print(response.text)

def get_pipeline(args, config):
    engine = get_engine(args, config)
    url = f"{engine['url']}v1/Pipelines/Get"
    payload = json.dumps({'name': args.name})
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, data=payload, headers=headers)
    print(response.text)

def delete_pipeline(args, config):
    engine = get_engine(args, config)
    url = f"{engine['url']}v1/Pipelines/Delete"
    payload = json.dumps({'name': args.name})
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, data=payload, headers=headers)
    if response.status_code == 200:
        print("Pipeline deleted successfully.")
    else:
        print("Failed to delete pipeline.")

def main():
    

    config = load_config()

    parser = argparse.ArgumentParser(description='Floom CLI (https://floom.ai)')
    subparsers = parser.add_subparsers()

    # Subparser for the 'apply' command
    parser_apply = subparsers.add_parser('apply', help='Apply yaml configuration.')
    parser_apply.add_argument('-f', '--file', help='File to apply')
    parser_apply.add_argument('-d', '--directory', help='Directory to apply')
    parser_apply.add_argument('-e', '--engine', help='Engine name')
    parser_apply.set_defaults(func=apply_file if '-f' in sys.argv or '--file' in sys.argv else apply_directory)

    # Subparser for the 'pipelines' command
    parser_pipelines = subparsers.add_parser('pipelines', help='Manage pipelines.')
    parser_pipelines_sub = parser_pipelines.add_subparsers()
    parser_pipelines_list = parser_pipelines_sub.add_parser('list', help='List pipelines.')
    parser_pipelines_list.add_argument('-e', '--engine', help='Engine name')
    parser_pipelines_list.set_defaults(func=list_pipelines)
    parser_pipelines_get = parser_pipelines_sub.add_parser('get', help='Get pipeline.')
    parser_pipelines_get.add_argument('-n', '--name', required=True, help='Pipeline name')
    parser_pipelines_get.set_defaults(func=get_pipeline)
    parser_pipelines_delete = parser_pipelines_sub.add_parser('delete', help='Delete pipeline.')
    parser_pipelines_delete.add_argument('-n', '--name', required=True, help='Pipeline name')
    parser_pipelines_delete.set_defaults(func=delete_pipeline)

    args = parser.parse_args()

    if hasattr(args, 'func'):
        args.func(args, config)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
