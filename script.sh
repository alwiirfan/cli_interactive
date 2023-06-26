# to build image cli_interactive
sudo docker build -t alwi09/cli-image .

# run and create container to cli_interactive and create network
sudo docker container run --name cli-container --network=host alwi09/cli-image

# 
sudo docker container logs cli-container

# look specification
sudo docker container inspect cli-container

## execute
sudo docker container exec -it container-name /bin/bash