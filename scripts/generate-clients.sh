#!/bin/bash

# Script to generate REST clients from OpenAPI specification
# Usage: ./scripts/generate-clients.sh [go|python|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OPENAPI_SPEC="$PROJECT_ROOT/docs/swagger.yaml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

generate_go_client() {
    log_info "Generating Go client..."

    if ! command -v openapi-generator-cli &> /dev/null; then
        log_error "openapi-generator-cli not found. Please install it first:"
        log_error "  npm install -g @openapitools/openapi-generator-cli"
        exit 1
    fi

    # Clean previous Go client
    rm -rf "$PROJECT_ROOT/clients/go"

    # Generate Go client using openapi-generator-cli
    openapi-generator-cli generate \
        -i "$OPENAPI_SPEC" \
        -g go \
        -o "$PROJECT_ROOT/clients/go" \
        --package-name ruleengine \
        --additional-properties=packageVersion=1.0.0,projectName=rule-engine-go-client,isGoSubmodule=true

    log_info "Go client generated successfully in clients/go/"
}

generate_python_client() {
    log_info "Generating Python client..."

    if ! command -v openapi-generator-cli &> /dev/null; then
        log_error "openapi-generator-cli not found. Please install it first:"
        log_error "  npm install -g @openapitools/openapi-generator-cli"
        exit 1
    fi

    # Clean previous Python client
    rm -rf "$PROJECT_ROOT/clients/python"

    # Generate Python client
    openapi-generator-cli generate \
        -i "$OPENAPI_SPEC" \
        -g python \
        -o "$PROJECT_ROOT/clients/python" \
        --package-name rule_engine_client \
        --additional-properties=packageVersion=1.0.0,projectName=rule-engine-client

    log_info "Python client generated successfully in clients/python/"
}

show_usage() {
    echo "Usage: $0 [go|python|all]"
    echo ""
    echo "Commands:"
    echo "  go      Generate Go client only"
    echo "  python  Generate Python client only"
    echo "  all     Generate both Go and Python clients"
    echo ""
    echo "Examples:"
    echo "  $0 go      # Generate Go client"
    echo "  $0 python  # Generate Python client"
    echo "  $0 all     # Generate both clients"
}

# Main logic
case "${1:-all}" in
    "go")
        generate_go_client
        ;;
    "python")
        generate_python_client
        ;;
    "all")
        generate_go_client
        generate_python_client
        ;;
    *)
        log_error "Invalid argument: $1"
        echo ""
        show_usage
        exit 1
        ;;
esac

log_info "Client generation completed!"