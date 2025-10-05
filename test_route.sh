#!/bin/bash

# 🎧 Audio Guide - Complete Workflow Test
# This script demonstrates the full route generation process

set -e

API_BASE="http://localhost:8080/api"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  🎧 Audio Guide Workflow Test${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Step 1: Check server health
echo -e "${YELLOW}📡 Step 1: Checking server health...${NC}"
HEALTH=$(curl -s ${API_BASE}/health)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Server is running${NC}"
    echo "$HEALTH" | python3 -m json.tool
else
    echo -e "❌ Server is not running! Start it with:"
    echo "   cd backend && go run cmd/server/main.go"
    exit 1
fi
echo ""

# Step 2: List available POIs
echo -e "${YELLOW}📍 Step 2: Checking available POIs...${NC}"
POI_COUNT=$(curl -s ${API_BASE}/pois | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data))" 2>/dev/null || echo "0")
echo -e "${GREEN}✅ Found $POI_COUNT POIs in database${NC}"

# Show some POIs
echo "Sample POIs:"
curl -s "${API_BASE}/pois?epoch=soviet" | python3 -m json.tool | head -50
echo ""

# Step 3: Generate a route
echo -e "${YELLOW}🚀 Step 3: Generating Soviet Architecture route...${NC}"
echo "Start point: Red Square (55.7539, 37.6208)"
echo "Duration: 60 minutes"
echo "Preferences: Soviet era, Architecture"
echo ""

ROUTE_JSON=$(curl -s -X POST ${API_BASE}/routes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "start_point": {
      "lat": 55.7539,
      "lon": 37.6208
    },
    "duration_minutes": 60,
    "epochs": ["soviet"],
    "interests": ["architecture"],
    "max_waypoints": 4
  }')

if [ $? -ne 0 ]; then
    echo -e "❌ Failed to generate route"
    exit 1
fi

# Check if route was created
ROUTE_ID=$(echo "$ROUTE_JSON" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('route_id', ''))" 2>/dev/null)

if [ -z "$ROUTE_ID" ]; then
    echo -e "❌ Error: No route_id in response"
    echo "$ROUTE_JSON" | python3 -m json.tool
    exit 1
fi

echo -e "${GREEN}✅ Route generated successfully!${NC}"
echo "Route ID: $ROUTE_ID"
echo ""
echo "Route details:"
echo "$ROUTE_JSON" | python3 -m json.tool
echo ""

# Step 4: Wait for content generation
echo -e "${YELLOW}⏳ Step 4: Waiting for content generation...${NC}"
echo "Content (text + audio) is being generated asynchronously."
echo "This takes about 30-60 seconds per waypoint."
echo ""

for i in {10..1}; do
    echo -ne "Checking in $i seconds...\r"
    sleep 1
done
echo ""

# Step 5: Check route status
echo -e "${YELLOW}📥 Step 5: Fetching complete route with content...${NC}"
COMPLETE_ROUTE=$(curl -s ${API_BASE}/routes/${ROUTE_ID})
echo "$COMPLETE_ROUTE" | python3 -m json.tool

# Check if content is ready
HAS_CONTENT=$(echo "$COMPLETE_ROUTE" | python3 -c "import sys, json; data=json.load(sys.stdin); wp=data.get('waypoints', []); print('yes' if len(wp) > 0 and wp[0].get('content') else 'no')" 2>/dev/null)

if [ "$HAS_CONTENT" = "yes" ]; then
    echo ""
    echo -e "${GREEN}✅ Content generation complete!${NC}"
    
    # Step 6: Try to download audio
    echo ""
    echo -e "${YELLOW}🔊 Step 6: Downloading audio files...${NC}"
    
    # Get first waypoint ID
    FIRST_WAYPOINT=$(echo "$COMPLETE_ROUTE" | python3 -c "import sys, json; data=json.load(sys.stdin); wp=data.get('waypoints', []); print(wp[0]['id'] if len(wp) > 0 else '')" 2>/dev/null)
    
    if [ -n "$FIRST_WAYPOINT" ]; then
        echo "Downloading audio for first waypoint: $FIRST_WAYPOINT"
        curl -s ${API_BASE}/audio/${FIRST_WAYPOINT} -o waypoint_audio.mp3
        
        if [ -f waypoint_audio.mp3 ] && [ -s waypoint_audio.mp3 ]; then
            echo -e "${GREEN}✅ Audio downloaded: waypoint_audio.mp3${NC}"
            echo "File size: $(ls -lh waypoint_audio.mp3 | awk '{print $5}')"
        else
            echo -e "${YELLOW}⚠️  Audio not ready yet. Try again in a minute.${NC}"
        fi
    fi
    
    # Step 7: Download complete route audio guide
    echo ""
    echo -e "${YELLOW}🎧 Step 7: Downloading complete route audio guide...${NC}"
    echo "Downloading full audio guide for route: $ROUTE_ID"
    
    curl -s ${API_BASE}/routes/${ROUTE_ID}/audio -o route_audio_guide.mp3
    
    if [ -f route_audio_guide.mp3 ] && [ -s route_audio_guide.mp3 ]; then
        echo -e "${GREEN}✅ Complete audio guide downloaded: route_audio_guide.mp3${NC}"
        echo "File size: $(ls -lh route_audio_guide.mp3 | awk '{print $5}')"
        echo ""
        echo "🎵 You can now play the audio guide:"
        echo "   afplay route_audio_guide.mp3  # macOS"
        echo "   mpg123 route_audio_guide.mp3  # Linux"
    else
        echo -e "${YELLOW}⚠️  Audio guide not ready yet. Try again in a minute.${NC}"
        echo "Check status with:"
        echo "  curl ${API_BASE}/routes/${ROUTE_ID}/audio"
    fi
else
    echo ""
    echo -e "${YELLOW}⚠️  Content is still being generated.${NC}"
    echo "Wait 1-2 minutes and check the route again:"
    echo "curl ${API_BASE}/routes/${ROUTE_ID} | python3 -m json.tool"
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  ✅ Test Complete${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Route ID: $ROUTE_ID"
echo ""
echo "To check route status:"
echo "  curl ${API_BASE}/routes/${ROUTE_ID} | python3 -m json.tool"
echo ""
echo "To download single waypoint audio:"
echo "  curl ${API_BASE}/audio/<waypoint_id> -o audio.mp3"
echo ""
echo "To download complete route audio guide:"
echo "  curl ${API_BASE}/routes/${ROUTE_ID}/audio -o route_guide.mp3"
echo ""
