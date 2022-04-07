ARG TERRAFORM_VERSION=latest
FROM hashicorp/terraform:$TERRAFORM_VERSION AS terraform

FROM golang:1.17.8-alpine3.15
RUN apk --no-cache add make git bash curl

# Install terraform
COPY --from=terraform /bin/terraform /usr/local/bin/

# Install tfupdate
ENV TFUPDATE_VERSION 0.6.5
RUN curl -fsSL https://github.com/minamijoyo/tfupdate/releases/download/v${TFUPDATE_VERSION}/tfupdate_${TFUPDATE_VERSION}_linux_amd64.tar.gz \
  | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/tfupdate

# Install tfmigrate
ENV TFMIGRATE_VERSION 0.3.2
RUN curl -fsSL https://github.com/minamijoyo/tfmigrate/releases/download/v${TFMIGRATE_VERSION}/tfmigrate_${TFMIGRATE_VERSION}_linux_amd64.tar.gz \
  | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/tfmigrate

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make install

ENTRYPOINT ["./entrypoint.sh"]
