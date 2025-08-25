#!/bin/bash

# Pre-commit checks to prevent regression of remediation fixes
# Run this script before committing to catch common violations

set -e

echo "🔍 Running pre-commit validation checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track overall success
OVERALL_SUCCESS=true

#############################################
# Check 1: No time.Now() in core logic
#############################################
echo -e "\n📅 Checking for time.Now() usage in core logic..."

CORE_PACKAGES=(
    "internal/game_state"
    "internal/ai"
    "internal/sprites"
    "internal/map_state"
    "internal/party_state"
)

TIME_NOW_VIOLATIONS=0
for package in "${CORE_PACKAGES[@]}"; do
    if [ -d "$package" ]; then
        violations=$(grep -r "time\.Now()" "$package" 2>/dev/null | grep -v "_test.go" | grep -v "// ALLOWED:" || true)
        if [ -n "$violations" ]; then
            echo -e "${RED}❌ Found time.Now() usage in $package:${NC}"
            echo "$violations"
            TIME_NOW_VIOLATIONS=1
        fi
    fi
done

if [ $TIME_NOW_VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ No time.Now() violations in core logic${NC}"
else
    echo -e "${RED}❌ time.Now() usage found in core logic packages${NC}"
    echo -e "${YELLOW}💡 Use GameState.DateTime or central GameClock instead${NC}"
    OVERALL_SUCCESS=false
fi

#############################################
# Check 2: log.Fatal has comments
#############################################
echo -e "\n🚨 Checking log.Fatal comments..."

FATAL_VIOLATIONS=0
fatal_calls=$(grep -rn "log\.Fatal" --include="*.go" . 2>/dev/null | grep -v "_test.go" || true)

if [ -n "$fatal_calls" ]; then
    while IFS= read -r line; do
        file_line=$(echo "$line" | cut -d: -f1,2)
        
        # Get the line before the log.Fatal call to check for comment
        file=$(echo "$line" | cut -d: -f1)
        line_num=$(echo "$line" | cut -d: -f2)
        prev_line_num=$((line_num - 1))
        
        if [ $prev_line_num -gt 0 ]; then
            prev_line=$(sed -n "${prev_line_num}p" "$file" 2>/dev/null || echo "")
            # Check if previous line has TODO or explanatory comment
            if echo "$prev_line" | grep -q "// TODO\|// .*:" && ! echo "$prev_line" | grep -q "^[[:space:]]*$"; then
                # Has appropriate comment
                continue
            fi
        fi
        
        echo -e "${RED}❌ log.Fatal without appropriate comment: $file_line${NC}"
        FATAL_VIOLATIONS=1
    done <<< "$fatal_calls"
fi

if [ $FATAL_VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ All log.Fatal calls have appropriate comments${NC}"
else
    echo -e "${RED}❌ Found log.Fatal calls without comments${NC}"
    echo -e "${YELLOW}💡 Add explanatory comment or '// TODO: soften to recoverable error'${NC}"
    OVERALL_SUCCESS=false
fi

#############################################
# Check 3: Import organization
#############################################
echo -e "\n📦 Checking import organization..."

# Check if goimports would make changes
IMPORT_VIOLATIONS=0
unformatted_files=$(goimports -l . 2>/dev/null | grep -v vendor/ || true)

if [ -n "$unformatted_files" ]; then
    echo -e "${RED}❌ Files with incorrect import formatting:${NC}"
    echo "$unformatted_files"
    echo -e "${YELLOW}💡 Run 'goimports -w .' to fix formatting${NC}"
    IMPORT_VIOLATIONS=1
    OVERALL_SUCCESS=false
else
    echo -e "${GREEN}✅ Import organization correct${NC}"
fi

#############################################
# Check 4: Package boundary violations
#############################################
echo -e "\n🏗️ Checking package boundary violations..."

BOUNDARY_VIOLATIONS=0
boundary_violations=$(grep -r "github.com/hajimehoshi/ebiten/v2" internal/game_state/ internal/ai/ internal/map_state/ 2>/dev/null | grep -v "_test.go" | grep -v "// ALLOWED:" || true)

if [ -n "$boundary_violations" ]; then
    echo -e "${RED}❌ Core logic packages importing Ebitengine:${NC}"
    echo "$boundary_violations"
    echo -e "${YELLOW}💡 Use DisplayManager or dependency injection instead${NC}"
    BOUNDARY_VIOLATIONS=1
    OVERALL_SUCCESS=false
else
    echo -e "${GREEN}✅ No package boundary violations${NC}"
fi

#############################################
# Check 5: Unnecessary import aliases
#############################################
echo -e "\n🔗 Checking for unnecessary import aliases..."

ALIAS_VIOLATIONS=0
bad_aliases=$(grep -r "references2\|ucolor\|mainscreen2" --include="*.go" . | grep "import" || true)

if [ -n "$bad_aliases" ]; then
    echo -e "${RED}❌ Found unnecessary import aliases:${NC}"
    echo "$bad_aliases"
    echo -e "${YELLOW}💡 Remove unnecessary aliases like references2, ucolor, mainscreen2${NC}"
    ALIAS_VIOLATIONS=1
    OVERALL_SUCCESS=false
else
    echo -e "${GREEN}✅ No unnecessary import aliases found${NC}"
fi

#############################################
# Check 6: Basic Go checks
#############################################
echo -e "\n🔧 Running basic Go checks..."

echo "Running go vet..."
if ! go vet ./...; then
    echo -e "${RED}❌ go vet failed${NC}"
    OVERALL_SUCCESS=false
else
    echo -e "${GREEN}✅ go vet passed${NC}"
fi

echo "Running go build..."
if ! go build ./...; then
    echo -e "${RED}❌ go build failed${NC}"
    OVERALL_SUCCESS=false
else
    echo -e "${GREEN}✅ go build passed${NC}"
fi

#############################################
# Summary
#############################################
echo -e "\n📋 Pre-commit Check Summary:"
echo "============================================"

if [ "$OVERALL_SUCCESS" = true ]; then
    echo -e "${GREEN}🎉 All checks passed! Ready to commit.${NC}"
    exit 0
else
    echo -e "${RED}❌ Some checks failed. Please fix the issues above before committing.${NC}"
    echo -e "\n🔧 Quick fixes:"
    echo "  • Run 'goimports -w .' for import formatting"
    echo "  • Add comments to log.Fatal calls"  
    echo "  • Replace time.Now() with GameState time sources"
    echo "  • Use DisplayManager instead of direct Ebitengine imports"
    echo "  • Remove unnecessary import aliases"
    echo -e "\n📚 See docs/CODE_REVIEW_CHECKLIST.md for details"
    exit 1
fi