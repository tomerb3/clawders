Task: Add a new Leveraged ETF Daily Opportunity Scanner page at route `/etf`

Context:
You are updating our existing Finviz-like web application running locally on port 8782. We need a dedicated page (`http://localhost:8782/etf`) tailored for daily check-ins to spot high-conviction pullbacks and oversold conditions in prominent leveraged ETFs.

Requirements & Specifications:

1. Target Route & Layout
- Route: `/etf`
- Integrate navigation links to `/etf` in the main header/navbar alongside existing pages.
- UI Theme: Match the existing dark/light Finviz-like aesthetic and table styling.
- Design for a "Daily Glance" workflow: clear visual indicators, sorted by best opportunity score by default.

2. ETF Universe (Leveraged Bull Funds)
Include prominent 2x and 3x leveraged ETFs across major indexes, tech, semiconductor, and sector funds, such as:
- Tech & Broad Market: TQQQ, UPRO, SOXL, TECL, USD, FNGU, BULZ
- Small Cap & Financials: TNA, FAS, DPST
- Biotech & Energy: LABU, ERX, GUSH, NUGT
(Make this universe easily expandable in a central config or JSON file).

3. Data Metrics & Dip Analysis Engine
For each ETF in the list, fetch or compute:
- Current Price & Today's % Change
- 1-Month Benchmark: % drop/gain relative to 1-Month High and 1-Month Average
- 3-Month Benchmark: % drop/gain relative to 3-Month High and 3-Month Average
- Technical Indicators: 14-day RSI, 20-day SMA, 50-day SMA distance.
- Sentiment Indicator: Aggregate news/market sentiment score (or Fear & Greed / stock sentiment metric if available in our data layer, otherwise derive a technical sentiment score from momentum + volume).

4. Signal & Confidence Scoring Logic
Implement a scoring algorithm that outputs a "Buy Confidence Score" (0–100%) and a Signal Status:
- "Strong Buy / Deep Dip" (e.g. Confidence > 75%): Price is significantly down from 1M/3M highs, RSI < 35, and broader trend is bullish.
- "Watchlist / Dip Forming" (Confidence 50%–74%): Moderate pullback (down 10–20% from 1M/3M peak).
- "Neutral / Hold" (Confidence 30%–49%): Extended or consolidating.
- "Avoid / Overbought" (Confidence < 30%): Near 1M/3M high or RSI > 70.

5. UI Components & Dashboard Table
The `/etf` page should display:
- Metric Cards at top:
  * Most Oversold Leveraged ETF Today
  * Market Sentiment Overview (Bullish/Bearish/Neutral)
  * Top Signal of the Day
- Interactive Data Table with columns:
  1. Ticker & Name
  2. Current Price
  3. Today's %
  4. % Off 1M High
  5. % Off 3M High
  6. 14-Day RSI
  7. Sentiment Rating
  8. Buy Confidence Meter (visual progress bar or badge with exact %)
  9. Action Signal Badge (STRONG BUY, WATCH, NEUTRAL, AVOID)
- Filter & Sorting Controls:
  * Quick Filters: "All", "Deep Dips (>15% Drop)", "RSI Oversold (<35)", "High Confidence (>70%)"
  * Column sorting for all numerical fields.

6. Backend / API Execution
- Leverage existing API routes or data fetchers used by the main app.
- Implement efficient caching (e.g., 5–15 min cache or daily fetch) so loading `/etf` is instant.

Please review the current codebase architecture, identify where to place the route handler, views/components, and data processing scripts, and implement the changes step-by-step.
