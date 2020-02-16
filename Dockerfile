FROM golang:1.13

ENV LOG_LEVEL=info
ENV PG_DSN="postgres://%v@%v:5432/%v?sslmode=disable"
ENV PG_DSN_USER_NAME="marciodasilva"
ENV PG_DSN_USER_PASSWORD=""
ENV PG_DSN_HOST="127.0.0.1"
ENV PG_DSN_SCHEMA="plapi"
# Configure the repo url so we can configure our work directory:
ENV REPO_URL=github.com/prosline/pl_users_api

# Setup out $GOPATH
ENV GOPATH=/app

ENV APP_PATH=$GOPATH/src/$REPO_URL

# /app/src/github.com/prosline/pl_users_api/src

# Copy the entire source code from the current directory to $WORKPATH
ENV WORKPATH=$APP_PATH/src
COPY src $WORKPATH
WORKDIR $WORKPATH

RUN go build -o users_api .

# Expose port 8081 to the world:
EXPOSE 8081

CMD ["./users_api"]