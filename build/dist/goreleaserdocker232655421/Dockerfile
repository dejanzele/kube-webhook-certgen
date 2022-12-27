FROM alpine:3.16.2

RUN apk add --update --no-cache libc6-compat

COPY kube-webhook-certgen /kube-webhook-certgen

ARG USER=default
ENV HOME /home/$USER

# install sudo as root
RUN apk add --update sudo libc6-compat

# add new user
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME

ENTRYPOINT ["/kube-webhook-certgen"]
