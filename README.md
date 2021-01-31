Web page parsing and analyzing 

<h1>How to start?</h1>

<h3>Clone</h3>

`git clone git@github.com:Dmitriy-Opria/re_web_page_analyzer.git &&`</br>
`cd re_web_page_analyzer`

<h3>Create .env file</h3>

`touch cmd/.env &&` </br>
`echo LOG_LEVEL=5 >> cmd/.env &&`</br>
`echo RELEASE=test >> cmd/.env &&`</br>
`echo VCS_REF=test >> cmd/.env &&`</br>
`echo API_LISTENER=0.0.0.0:9088 >> cmd/.env &&`</br>
`echo API_WORKER_COUNT=50 >> cmd/.env`

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

<h1>Questions</h1>
1. It was not clear for me, how to check links for inaccessibility, 
I checked just for GET method. 
And I think it should be enough, maybe this logic should be changed

<h1>Improvements</h1>
<ul>
<li> App should works better with SPA.</li>
<li> Key words for checking login forms should be improved.</li>
<li> Client part should be moved to the different front-end application.</li>
<li> Html pages need to be prepared better.</li>
<li> Login checks should not check only forms.</li>
</ul>