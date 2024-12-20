# Step 1: Build the TypeScript application
FROM node:18-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json (or yarn.lock) files
COPY package*.json ./

# Install the dependencies
RUN npm install

# Copy the rest of the TypeScript source code
COPY . .

# Build the TypeScript application
RUN npm run build

# Step 2: Create a minimal image for running the application
FROM node:18-alpine

# Set the working directory inside the container
WORKDIR /app

# Install necessary production dependencies (if needed)
RUN npm install --production

# Copy the build output from the builder stage
COPY --from=builder /app/dist /app/dist

# Expose the port the application listens on
EXPOSE 3000

# Run the TypeScript application (the entry point is in the dist folder after building)
CMD ["node", "dist/server.js"]
