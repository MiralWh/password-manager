FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production
ENV BACKEND_URL=http://backend:8080
COPY --from=build /app/package.json    ./
COPY --from=build /app/node_modules    ./node_modules
COPY --from=build /app/dist            ./dist
COPY --from=build /app/vite.config.ts  ./
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
  CMD wget -qO- http://localhost:3000/ || exit 1
CMD ["npx", "vite", "preview", "--host", "0.0.0.0", "--port", "3000"]
