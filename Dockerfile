FROM node:9-alpine

# Create app directory
WORKDIR /usr/src/app

# Bundle app source
COPY . .

RUN npm install

EXPOSE 8081
CMD ["node", "app.js"]
