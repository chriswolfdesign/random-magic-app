FROM node:latest

COPY package.json package.json
COPY public public

RUN npm install
COPY src src

CMD ["npm", "run", "start"]
