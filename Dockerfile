ARG TERRAFORM_VERSION=latest
FROM hashicorp/terraform:$TERRAFORM_VERSION AS terraform

FROM golang:1.19-alpine3.17
RUN apk --no-cache add make git bash curl jq

# A workaround for a permission issue of git.
# Since UIDs are different between host and container,
# the .git directory is untrusted by default.
# We need to allow it explicitly.
# https://github.com/actions/checkout/issues/760
RUN git config --global --add safe.directory /work

# Install terraform
COPY --from=terraform /bin/terraform /usr/local/bin/

# Install tfupdate
ENV TFUPDATE_VERSION 0.6.5
RUN curl -fsSL https://github.com/minamijoyo/tfupdate/releases/download/v${TFUPDATE_VERSION}/tfupdate_${TFUPDATE_VERSION}_linux_amd64.tar.gz \
  | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/tfupdate

# Install tfmigrate
ENV TFMIGRATE_VERSION 0.3.3
RUN curl -fsSL https://github.com/minamijoyo/tfmigrate/releases/download/v${TFMIGRATE_VERSION}/tfmigrate_${TFMIGRATE_VERSION}_linux_amd64.tar.gz \
  | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/tfmigrate

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make install

ENTRYPOINT ["./entrypoint.sh"]
