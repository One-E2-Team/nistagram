FROM node as build
ENV NODE_ENV=development
COPY frontend /nistagram
COPY agent_frontend /agent
WORKDIR /nistagram
RUN npm install -d && npm run build
WORKDIR /agent
RUN npm install -d && npm run build


FROM nginx as final
COPY --from=build /nistagram/dist /usr/share/nginx/html/dist/web
EXPOSE 80
EXPOSE 443
STOPSIGNAL SIGTERM
CMD ["nginx", "-g", "daemon off;"]

FROM nginx as agent-final
COPY --from=build /agent/dist /usr/share/nginx/html/
EXPOSE 80
EXPOSE 443
STOPSIGNAL SIGTERM
CMD ["nginx", "-g", "daemon off;"]