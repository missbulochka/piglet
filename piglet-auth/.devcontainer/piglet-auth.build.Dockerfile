FROM golang:1.21.3 AS builder

RUN apt-get update -y && apt-get upgrade -y \
    	&& apt-get clean \
        && rm -rf /var/lib/apt/lists/*

ARG USERNAME=builder
ARG USER_UID=1000
ARG USER_GID=${USER_UID}

RUN groupadd --gid ${USER_GID} ${USERNAME} \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m -s /bin/bash ${USERNAME} \
    && chown -R ${USER_UID}:${USER_GID} /home/${USERNAME} \
    && mkdir -p /etc/sudoers.d/ \
    && echo ${USERNAME} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}

USER ${USERNAME}:${USERNAME}

WORKDIR /workspaces/dev_piglet
COPY ../go.mod .

RUN go mod download

EXPOSE 8080
