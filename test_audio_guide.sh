#!/bin/bash

# 🎧 Quick Audio Guide Test
# This script tests the new route audio guide endpoint

set -e

API_BASE="http://localhost:8080/api"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  🎧 Audio Guide Download Test${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if route ID is provided
if [ -z "$1" ]; then
    echo -e "${YELLOW}Usage: $0 <route_id>${NC}"
    echo ""
    echo "First, generate a route:"
    echo "  ./test_route.sh"
    echo ""
    echo "Then use this script to download the audio guide:"
    echo "  $0 <route_id>"
    echo ""
    exit 1
fi

ROUTE_ID=$1

echo -e "${YELLOW}📡 Checking server health...${NC}"
if ! curl -s ${API_BASE}/health > /dev/null; then
    echo -e "${RED}❌ Server is not running!${NC}"
    echo "Start it with: cd backend && go run cmd/server/main.go"
    exit 1
fi
echo -e "${GREEN}✅ Server is running${NC}"
echo ""

echo -e "${YELLOW}🔍 Checking route status...${NC}"
ROUTE_JSON=$(curl -s ${API_BASE}/routes/${ROUTE_ID})

# Check if route exists
if echo "$ROUTE_JSON" | grep -q "error"; then
    echo -e "${RED}❌ Route not found or error occurred${NC}"
    echo "$ROUTE_JSON" | python3 -m json.tool
    exit 1
fi

echo -e "${GREEN}✅ Route found${NC}"
echo ""

# Show route details
ROUTE_NAME=$(echo "$ROUTE_JSON" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('name', 'Unknown'))" 2>/dev/null)
WAYPOINT_COUNT=$(echo "$ROUTE_JSON" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data.get('waypoints', [])))" 2>/dev/null)

echo "Route: $ROUTE_NAME"
echo "Waypoints: $WAYPOINT_COUNT"
echo ""

echo -e "${YELLOW}🎧 Downloading complete audio guide...${NC}"
echo "Endpoint: ${API_BASE}/routes/${ROUTE_ID}/audio"
echo ""

OUTPUT_FILE="route_${ROUTE_ID:0:8}_guide.mp3"

# Download with verbose output to see response
HTTP_CODE=$(curl -s -w "%{http_code}" -o "$OUTPUT_FILE" ${API_BASE}/routes/${ROUTE_ID}/audio)

if [ "$HTTP_CODE" = "200" ]; then
    if [ -f "$OUTPUT_FILE" ] && [ -s "$OUTPUT_FILE" ]; then
        FILE_SIZE=$(ls -lh "$OUTPUT_FILE" | awk '{print $5}')
        echo -e "${GREEN}✅ Audio guide downloaded successfully!${NC}"
        echo ""
        echo "File: $OUTPUT_FILE"
        echo "Size: $FILE_SIZE"
        echo ""
        echo -e "${BLUE}🎵 Play the audio guide:${NC}"
        echo "  afplay $OUTPUT_FILE  # macOS"
        echo "  mpg123 $OUTPUT_FILE  # Linux"
        echo "  vlc $OUTPUT_FILE     # VLC player"
        echo ""
        
        # Try to play automatically on macOS
        if command -v afplay &> /dev/null; then
            echo -e "${YELLOW}Playing audio guide...${NC}"
            afplay "$OUTPUT_FILE"
        fi
    else
        echo -e "${RED}❌ Downloaded file is empty${NC}"
        rm -f "$OUTPUT_FILE"
        exit 1
    fi
elif [ "$HTTP_CODE" = "206" ]; then
    echo -e "${YELLOW}⚠️  Partial content (HTTP 206)${NC}"
    echo "Some audio files are still being generated."
    echo ""
    cat "$OUTPUT_FILE" | python3 -m json.tool
    rm -f "$OUTPUT_FILE"
    echo ""
    echo "Wait 1-2 minutes and try again:"
    echo "  $0 $ROUTE_ID"
elif [ "$HTTP_CODE" = "404" ]; then
    echo -e "${YELLOW}⚠️  Audio not ready yet (HTTP 404)${NC}"
    echo "Audio is still being generated."
    echo ""
    if [ -f "$OUTPUT_FILE" ]; then
        cat "$OUTPUT_FILE" | python3 -m json.tool
        rm -f "$OUTPUT_FILE"
    fi
    echo ""
    echo "Wait 1-2 minutes and try again:"
    echo "  $0 $ROUTE_ID"
else
    echo -e "${RED}❌ Error: HTTP $HTTP_CODE${NC}"
    if [ -f "$OUTPUT_FILE" ]; then
        cat "$OUTPUT_FILE"
        rm -f "$OUTPUT_FILE"
    fi
    exit 1
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  ✅ Test Complete${NC}"
echo -e "${BLUE}========================================${NC}"
