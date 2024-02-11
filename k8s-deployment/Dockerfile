FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Install Python and pip
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*


# Copy files from the first image
COPY --from=arearnest7/video-recog-oci:latest /app /app/deployment1

# Copy files from the second image
COPY --from=arearnest7/video-decoder-oci:latest /app /app/deployment2

# Copy files from the third image
COPY --from=arearnest7/video-streaming-oci:latest /app /app/deploymen3

COPY ./requirements.txt requirements.txt

RUN pip install -r requirements.txt

# Set the working directory
WORKDIR /app





CMD ["sleep 3600"]
