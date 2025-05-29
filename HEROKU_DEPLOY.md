# Heroku Deployment Guide

## Prerequisites
1. Install Heroku CLI: https://devcenter.heroku.com/articles/heroku-cli
2. Create a Heroku account
3. Get your API key from: https://dashboard.heroku.com/account

## Environment Variables Setup
Set these in your Heroku dashboard or via CLI:

```bash
# Replace YOUR_HEROKU_API_KEY with your actual API key
heroku config:set HEROKU_API_KEY=YOUR_HEROKU_API_KEY

# Application configuration
heroku config:set ENVIRONMENT=production
heroku config:set LOG_LEVEL=info
heroku config:set BASE_URL=https://your-app-name.herokuapp.com

# Database will be automatically set when you add Heroku Postgres
```

## Deployment Steps

### 1. Login to Heroku
```bash
heroku login
```

### 2. Create Heroku App
```bash
heroku create your-url-shortener-app
```

### 3. Add PostgreSQL Database
```bash
heroku addons:create heroku-postgresql:mini
```

### 4. Deploy Application
```bash
git add .
git commit -m "Deploy to Heroku"
git push heroku main
```

### 5. Check Application Status
```bash
heroku ps:scale web=1
heroku logs --tail
heroku open
```

## Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | 8080 | No (Heroku sets this) |
| `DATABASE_URL` | PostgreSQL connection string | - | Yes (Auto-set by Heroku) |
| `BASE_URL` | Base URL for short links | http://localhost:8080 | Yes |
| `ENVIRONMENT` | App environment | development | No |
| `LOG_LEVEL` | Logging level | info | No |

## Monitoring and Logs

### View Logs
```bash
heroku logs --tail
```

### Check App Status
```bash
heroku ps
```

### Open App
```bash
heroku open
```

## Scaling

### Scale Web Dynos
```bash
heroku ps:scale web=2
```

### Upgrade Database
```bash
heroku addons:upgrade heroku-postgresql:basic
```

## Troubleshooting

### Common Issues

1. **Build Failures**: Check that all dependencies are in go.mod
2. **Database Connection**: Ensure DATABASE_URL is set
3. **Port Issues**: Heroku sets PORT automatically, don't hardcode it

### Debug Commands
```bash
# Check environment variables
heroku config

# View recent logs
heroku logs --num=500

# Run one-off commands
heroku run bash
```
