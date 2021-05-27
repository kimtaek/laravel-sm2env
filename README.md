### AWS Secret manager to laravel .env

### usage

```bash
AWS_SECRET_ID=LARAVEL-ENV \
AWS_REGION=ap-northeast-2 \
AWS_ACCESS_KEY=AKIAY7V... \
AWS_SECRET_KEY=njtjtMRUET... \
./laravel-sm2env
```

##### [required] AWS_SECRET_ID (default: `LARAVEL-ENV`)
##### [required_if:AWS_PROFILE=""] AWS_ACCESS_KEY (default: ``)
##### [required_if:AWS_PROFILE=""] AWS_SECRET_KEY (default: ``)
##### [required_if:AWS_ACCESS_KEY=""] AWS_PROFILE (default: ``)
##### [optional] AWS_REGION (default: `us-west-2`)
##### [optional] AWS_SECRET_VERSION (default: `AWSCURRENT`)
##### [optional] FILEPATH (default: `.env`)
