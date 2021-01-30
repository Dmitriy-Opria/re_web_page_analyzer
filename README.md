Web page parsing and analyzing 

<h1>How to start?</h1>

<h3>Clone</h3>

`git clone git@github.com:Dmitriy-Opria/re_web_page_analyzer.git
cd re_web_page_analyzer`

<h3>Create .env file</h3>

`touch cmd/.env &&
echo LOG_LEVEL=5 >> cmd/.env &&
echo RELEASE=test >> cmd/.env &&
echo VCS_REF=test >> cmd/.env &&
echo API_LISTENER=0.0.0.0:9088 >> cmd/.env &&
echo API_WORKER_COUNT=50 >> cmd/.env`

<h3>Build docker image</h3>

`docker build -t re_web_page_analyzer -f Dockerfile .`

<h3>Run docker image</h3>

`docker run -it -d -p 9088:9088 --env-file=cmd/.env --name re_web_page_analyzer re_web_page_analyzer:latest`

<h3>Make test request</h3>

`curl --location --request POST 'http://0.0.0.0:9088/api/v1/parsing/page/analyze' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url": "https://www.w3schools.com/"
}'`
