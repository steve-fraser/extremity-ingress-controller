FROM golang:1.17-alpine AS build


ENV GO111MODULE=on
WORKDIR /src/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/extremity ./cmd


FROM haproxytech/kubernetes-ingress

RUN cd /tmp/ && wget https://github.com/projectcalico/bird/releases/download/v0.3.2/bird
RUN cd /tmp/ && wget https://github.com/projectcalico/bird/releases/download/v0.3.2/bird6
RUN cd /tmp/ && wget https://github.com/projectcalico/bird/releases/download/v0.3.2/birdcl

# Link binaries to standard places /usr/sbin and /usr/bin
RUN ln -s /tmp/bird /usr/sbin/bird
RUN ln -s /tmp/bird6 /usr/sbin/bird6
RUN ln -s /tmp/birdcl /usr/bin/birdcl

# Create dirs needed for BIRD runtime
RUN mkdir -p /etc/bird
RUN mkdir -p /etc/bird6
RUN mkdir -p /usr/local/var/run
RUN mkdir -p /usr/local/etc

# Copy in global BIRD and BIRD6 configs
ADD birdy/bird.conf /etc/bird.conf
ADD birdy/bird6.conf /etc/bird6.conf
RUN ln -s /etc/bird.conf /usr/local/etc/bird.conf
RUN ln -s /etc/bird6.conf /usr/local/etc/bird6.conf

RUN chmod +x /tmp/bird
RUN chmod +x /tmp/bird6
RUN chmod +x /tmp/birdcl

COPY --from=build /bin/extremity /tmp/extremity

RUN ln -s /tmp/extremity /usr/bin/extremity
RUN chmod +x /tmp/extremity


RUN sed -i '/set -e/a bird' start.sh
RUN sed -i '/bird/a nohup extremity &' start.sh

RUN apk update && apk add curl git

RUN mkdir /tmp/manifests
ADD manifests/* /tmp/manifests
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.15.1/bin/linux/amd64/kubectl
RUN chmod u+x kubectl && mv kubectl /usr/bin/kubectl

