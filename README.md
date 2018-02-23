# GitHub Webhook


```
docker build -t gcr.io/hightowerlabs/github-webhook:v1 .
```

```
gcloud docker -- push gcr.io/hightowerlabs/github-webhook:v1
```

```
docker tag gcr.io/hightowerlabs/github-webhook:v1 \
  registry.heroku.com/github-webhook-heroku/web
```

```
docker push registry.heroku.com/github-webhook-heroku/web
```
