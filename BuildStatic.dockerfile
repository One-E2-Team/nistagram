FROM node as build
ENV NODE_ENV=development
COPY frontend /app
WORKDIR /app
RUN npm install && npm run build


FROM nginx as final
COPY --from=build /app/dist /usr/share/nginx/html/dist
EXPOSE 80
STOPSIGNAL SIGTERM
CMD ["nginx", "-g", "daemon off;"]