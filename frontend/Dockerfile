# frontend/Dockerfile

# Этап 1: Сборка зависимостей
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* yarn.lock* ./
# Если используете npm:
RUN npm ci
# Если yarn: RUN yarn install --frozen-lockfile

# Этап 2: Сборка приложения
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Отключаем телеметрию Next.js
ENV NEXT_TELEMETRY_DISABLED 1

# Собираем проект
RUN npm run build

# Этап 3: Запуск (Runner)
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

# Создаем пользователя, чтобы не запускать от root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Копируем необходимые файлы из билдера
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000
# Важно: 0.0.0.0 чтобы контейнер был доступен извне
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]