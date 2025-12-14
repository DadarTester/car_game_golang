FROM golang:1.24-bookworm AS builder

RUN apt-get update && apt-get install -y \
    libgl1-mesa-dev \
    libx11-dev \
    libx11-xcb-dev \
    libxcb-randr0-dev \
    libxcb-xfixes0-dev \
    libxcb-xinput-dev \
    libxcb-shape0-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxext-dev \
    libxxf86vm-dev \
    libasound2-dev \
    libpulse-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -v -ldflags='-w -s' -o main .

FROM ubuntu:22.04

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libgl1-mesa-glx \
    libgl1 \
    libglx0 \
    libgles2 \
    libegl1 \
    libglu1-mesa \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcb-randr0 \
    libxcb-xfixes0 \
    libxcb-shape0 \
    libxcb-xinput0 \
    libxcursor1 \
    libxi6 \
    libxinerama1 \
    libxrandr2 \
    libxext6 \
    libxxf86vm1 \
    libasound2 \
    libpulse0 \
    libwayland-client0 \
    libwayland-cursor0 \
    libwayland-egl1 \
    libxkbcommon0 \
    xvfb \
    mesa-utils \
    && rm -rf /var/lib/apt/lists/* \
    && ldconfig


RUN useradd -m -u 1000 appuser
USER appuser
WORKDIR /home/appuser

COPY --from=builder --chown=appuser:appuser /app/main /home/appuser/main

CMD ["./main"]