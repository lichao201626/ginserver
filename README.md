#### docker command
sudo service docker start
docker build -t server:1 .
docker run -d -t 8888:8888 server:1


# quanta_lab_aip
Quanta Lab Automatic Investment Plan

# install dependency
make deps
# build binary file
make build
# run test cases
make tests
# run app
make run
