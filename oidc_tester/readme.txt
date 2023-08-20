build:
 docker build --tag oidctester:latest .


docker run -p 8000:8000   -v /Users/avinashr/code/golang/oidc_tester/src/conf/:/oidc_tester/src/conf/  oidctester:latest
docker run -p 9000:9000   -v /Users/avinashr/code/golang/oidc_tester/src/conf/:/oidc_tester/src/conf/  oidctester:latest