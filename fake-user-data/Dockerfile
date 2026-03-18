FROM node:18-alpine

WORKDIR /app

COPY package.json ./
RUN npm install

COPY generate.js ./

# Data will be written to /data volume
VOLUME /data

CMD ["node", "generate.js"]
