# 使用官方Node镜像作为构建环境
FROM node:16-alpine AS build

WORKDIR /app

# 复制package.json和package-lock.json/yarn.lock
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制源码
COPY . .

# 构建前端应用
RUN npm run build

# 使用nginx镜像来部署前端
FROM nginx:alpine

# 复制构建产物到nginx的html目录
COPY --from=build /app/dist /usr/share/nginx/html

# 默认端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]