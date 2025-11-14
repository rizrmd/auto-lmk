# LLM Provider Research for Auto LMK

**Date:** 2025-11-14
**Task:** Day 2-3 - LLM Provider Research
**Status:** Pending Implementation

## Objective

Select the optimal LLM provider for our WhatsApp bot that handles car sales inquiries in **Bahasa Indonesia**.

## Requirements

### Primary Requirements

1. **Bahasa Indonesia Support**
   - Native or good quality Bahasa Indonesia understanding
   - Ability to generate natural Bahasa Indonesia responses
   - Understanding of automotive terminology in Indonesian (matic, manual, OTR, CBU, CKD, etc.)

2. **Function Calling**
   - Support for function/tool calling
   - Ability to call backend APIs (searchCars, getCarDetails, listInventory, createLead)
   - Return structured data for function parameters

3. **Cost Efficiency**
   - Target: Reasonable cost per conversation
   - Estimate: ~10-20 conversations per day per tenant
   - Estimate: ~5-10 messages per conversation
   - Need to calculate cost per month per tenant

4. **Latency**
   - Target: < 3-5 seconds response time for WhatsApp
   - Acceptable for asynchronous chat (not real-time)

5. **Context Window**
   - Need to maintain conversation history
   - Estimate: ~5-10 message pairs per session
   - Include car inventory context (variable size)

## Providers to Test

### 1. OpenAI (Recommended in Roadmap)

**Model Options:**
- GPT-4o (multimodal, latest)
- GPT-4-turbo
- GPT-3.5-turbo (cheaper alternative)

**Pros:**
- Industry standard
- Excellent function calling support
- Good multilingual capabilities
- Well-documented

**Cons:**
- Potentially higher cost
- Requires API key

**Testing Required:**
- [ ] Test Bahasa Indonesia quality
- [ ] Test function calling with sample car queries
- [ ] Calculate cost per 1M tokens
- [ ] Measure average response latency
- [ ] Test with Indonesian automotive terms

### 2. Anthropic Claude

**Model Options:**
- Claude 3.5 Sonnet
- Claude 3 Opus
- Claude 3 Haiku (cheaper, faster)

**Pros:**
- Excellent reasoning capabilities
- Good safety features
- Competitive pricing (especially Haiku)
- Strong function calling (tool use)

**Cons:**
- May have accessibility issues in Indonesia
- Newer, less established in Indonesian market

**Testing Required:**
- [ ] Verify API accessibility from Indonesia
- [ ] Test Bahasa Indonesia quality
- [ ] Test tool use with car search functions
- [ ] Calculate cost per 1M tokens
- [ ] Compare quality vs OpenAI

### 3. Local/Open Source Models (Future Consideration)

**Options:**
- LLaMA 3 (Meta)
- Mistral
- Indonesian-specific fine-tuned models

**Pros:**
- Lower long-term costs
- Data privacy (self-hosted)
- No API rate limits

**Cons:**
- Requires infrastructure (GPU)
- More complex setup
- May need fine-tuning for Indonesian
- Function calling support varies

**Decision:** Defer to Phase 2 unless cloud costs become prohibitive

## Test Scenarios

### Scenario 1: Customer Inquiry (Natural Language)
**Input:** "Ada mobil Toyota budget 200 juta ga?"
**Expected:**
- Understand intent: search for Toyota cars under 200 million
- Call searchCars function with filters: {brand: "Toyota", max_price: 200000000}
- Return natural response with results

### Scenario 2: Sales Query (Internal)
**Input:** "List semua mobil yang harga di bawah 300jt"
**Expected:**
- Identify as sales user (more access)
- Call searchCars with filter: {max_price: 300000000}
- Return complete inventory list

### Scenario 3: Specific Car Details
**Input:** "Yang Avanza 2023 ada foto nya ga?"
**Expected:**
- Understand: request for Avanza 2023 with photos
- Call searchCars or getCarDetails
- Respond with availability and indicate photos will be sent

### Scenario 4: Automotive Terminology
**Input:** "Cari mobil matic, bensin, tahun 2020 ke atas"
**Expected:**
- Understand "matic" = automatic transmission
- Parse year filter: >= 2020
- Call searchCars with {transmission: "automatic", fuel_type: "bensin", min_year: 2020}

### Scenario 5: Multi-turn Conversation
**Turn 1:** "Ada mobil apa aja?"
**Turn 2:** "Yang SUV aja"
**Turn 3:** "Yang paling murah berapa?"
**Expected:**
- Maintain context across turns
- Progressive filtering
- Remember previous queries

## Function Definitions for Testing

```json
{
  "name": "searchCars",
  "description": "Search available cars with optional filters",
  "parameters": {
    "type": "object",
    "properties": {
      "brand": {"type": "string", "description": "Car brand (Toyota, Honda, etc)"},
      "model": {"type": "string", "description": "Car model (Avanza, CR-V, etc)"},
      "max_price": {"type": "integer", "description": "Maximum price in IDR"},
      "min_year": {"type": "integer", "description": "Minimum year"},
      "transmission": {"type": "string", "enum": ["manual", "automatic"], "description": "Transmission type"},
      "fuel_type": {"type": "string", "enum": ["bensin", "diesel", "hybrid", "electric"], "description": "Fuel type"}
    }
  }
}
```

```json
{
  "name": "getCarDetails",
  "description": "Get detailed information about a specific car",
  "parameters": {
    "type": "object",
    "properties": {
      "car_id": {"type": "integer", "description": "Car ID"}
    },
    "required": ["car_id"]
  }
}
```

```json
{
  "name": "createLead",
  "description": "Create a lead when customer shows interest",
  "parameters": {
    "type": "object",
    "properties": {
      "phone_number": {"type": "string"},
      "name": {"type": "string"},
      "interested_car_id": {"type": "integer"}
    },
    "required": ["phone_number"]
  }
}
```

## Cost Estimation

### Assumptions:
- 10 tenants
- 15 conversations/day per tenant
- 7 messages per conversation (avg)
- 200 tokens per message (input + output)
- 30 days/month

**Total monthly tokens:** 10 × 15 × 7 × 200 × 30 = 6,300,000 tokens (~6.3M tokens/month)

### OpenAI Pricing (as of 2024):
- GPT-4o: $5/1M input, $15/1M output
- GPT-3.5-turbo: $0.50/1M input, $1.50/1M output

**Estimated Monthly Cost:**
- GPT-4o: ~$63 (input) + ~$189 (output) = **~$252/month**
- GPT-3.5-turbo: ~$3.15 + ~$9.45 = **~$12.60/month**

### Anthropic Pricing:
- Claude 3.5 Sonnet: $3/1M input, $15/1M output = **~$113/month**
- Claude 3 Haiku: $0.25/1M input, $1.25/1M output = **~$9.45/month**

## Decision Criteria

| Criteria | Weight | OpenAI GPT-4o | OpenAI GPT-3.5 | Claude 3.5 Sonnet | Claude 3 Haiku |
|----------|--------|---------------|----------------|-------------------|----------------|
| Indonesian Quality | 30% | ? | ? | ? | ? |
| Function Calling | 25% | Excellent | Good | Excellent | Good |
| Cost | 20% | Low | High | Medium | High |
| Latency | 15% | ? | ? | ? | ? |
| Availability (ID) | 10% | Good | Good | ? | ? |
| **Total** | 100% | TBD | TBD | TBD | TBD |

## Testing Plan

### Phase 1: Quick Validation (Day 2)
1. Sign up for OpenAI API
2. Sign up for Anthropic API (if accessible)
3. Run 10 test prompts in Bahasa Indonesia for each model
4. Test basic function calling
5. Measure response times

### Phase 2: Detailed Testing (Day 3)
1. Test all 5 scenarios above
2. Calculate actual token usage
3. Evaluate response quality (scale 1-10)
4. Test edge cases
5. Document findings

### Phase 3: Decision (End of Day 3)
1. Score each provider based on criteria
2. Make final selection
3. Document decision in this file
4. Update .env.example with chosen provider

## Next Steps

- [ ] Obtain OpenAI API key
- [ ] Obtain Anthropic API key (if accessible)
- [ ] Create test script for evaluation
- [ ] Run all test scenarios
- [ ] Calculate actual costs based on results
- [ ] Make final decision
- [ ] Update configuration

## Decision Log

**Status:** Pending

**Final Decision:** TBD

**Selected Provider:** TBD

**Rationale:** TBD

**Date Decided:** TBD

---

**Note:** This research is critical for Week 1 completion and Week 4-5 bot development.
