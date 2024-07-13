# Problem Statement 

#### General-purpose HTTP load-testing and benchmarking library

Write a program that serves as a general-purpose HTTP load-testing and benchmarking library. This is an open-ended project.

#### Bare minimum requirements:
- Takes an HTTP address as input
- Supports a --qps flag to generate requests at a given fixed QPS
- Reports latencies and error rates
- Makes a buildable Docker image
#### Restrictions:
- Takes an HTTP address as input
- Supports a --qps flag to generate requests at a given fixed QPS
- Reports latencies and error rates
- Makes a buildable Docker image
#### Evaluation criteria:
- Thoughtful design of what features such a library should provide. You are expected to implement additional features and functionality beyond the bare minimum requirements.
- Must provide accurate results. Your results will be compared against a reference implementation.
- API design
- Code, test, and documentation quality

## Design

### Language choice

**Go** stands out as the optimal language for developing a general-purpose HTTP load-testing and benchmarking library. It offers a good balance of performance, ease of use, and concurrency support, making it well-suited for handling high loads and providing accurate benchmarking results.

Factors considered for comparison:
| Category        | Go      | C++        | Python    | TypeScript |
|-----------------|---------|------------|-----------|------------|
| Performance     | High    | Very High  | Moderate  | Moderate   |
| Concurrency     | Excellent | Good     | Moderate  | Good       |
| Ease of Use     | High    | Low        | Very High | High       |
| Ecosystem       | Growing | Mature     | Very Mature | Mature    |
| Development Speed | Fast  | Slow       | Very Fast | Fast       |

Another factor is Developer Experience (high imapact) - My Strongest language is Java, I am equally skilled in other langauges Python, Go & TypeScript.

### Core Features
1. Input Handling: Accept an HTTP address as input.
2. QPS Control: Support a --qps flag to generate requests at a given fixed QPS (Queries Per Second).
3. Metrics Reporting: Report latencies and error rates.
4. Concurrency: Handle multiple concurrent requests.
5. Dockerization: Provide a Dockerfile for easy deployment.
### Additional Features
1. Configurable Request Headers: Allow users to specify custom headers.
2. Request Methods: Support different HTTP methods (GET, POST, etc.).
3. Payload support: passing request body especially usefully for DB write load testing, LLM performace testing etc.
4. Output Formats: Provide options for different output formats (JSON, plain text).
5. Graceful Shutdown: Handle interruptions gracefully and report partial results.

### Advanced Features
1. [File Payloads](#file-payload-feature-prioritization-analysis)

### API Design
- **CLI Arguments**: Use flags for configuration (e.g., --url, --qps, --method, --headers).

- **Metrics**: Collect and display metrics such as average latency, 95th percentile latency, and error rates.






<br/><br/>

### TLDR! Additiona Documentation
#### File Payload feature Prioritization analysis

categorizing the feature of files as payloads in an HTTP load-testing and benchmarking library using product management frameworks, such as RICE, Kano Model, and others. 

#### **RICE Scoring Model**
The RICE scoring model evaluates features based on Reach, Impact, Confidence, and Effort. Let's apply this to the file upload feature:
- Reach: How many users will benefit from the file upload feature? Given that many applications require file uploads (e.g., e-commerce, cloud storage), the reach is high.
    - Score: 8/10
- Impact: How significant is the impact of this feature on user satisfaction and application performance? File uploads are crucial for many use cases, significantly enhancing the tool's utility.
    - Score: 9/10
- Confidence: How confident are we in our estimates for reach and impact? Given the widespread need for file uploads and existing data on its importance, confidence is high.
    - Score: 8/10
- Effort: How much effort is required to implement this feature? Implementing file uploads, especially handling various file types and sizes, requires moderate to high effort.
    - Score: 4/10
  
**RICE Score Calculation**:

RICE Score = Reach × Impact × Confidence / Effort
           = 8 × 9 × 8 / 4
           = 144 / 4
           = 36

RICE Score = 36

#### **Kano Model**
The Kano Model categorizes features based on their ability to satisfy customers:
- Basic (Threshold) Features: These are essential features that users expect. For a load-testing tool, basic payload handling might be expected, but file uploads could be seen as an advanced feature.
    - File Uploads: Not a basic feature but expected in advanced tools.
- Performance Features: These features provide a proportional increase in satisfaction. File uploads can significantly enhance the tool's performance and user satisfaction by enabling comprehensive testing scenarios.
    - File Uploads: Performance feature.
- Excitement Features: These features delight users and exceed their expectations. For some users, the ability to handle complex file uploads might be an excitement feature.
    - File Uploads: Could be an excitement feature for advanced users.

#### **MoSCoW Method**
The MoSCoW method categorizes features into Must-Have, Should-Have, Could-Have, and Won't-Have:
- Must-Have: Essential for the product's basic functionality.
    - File Uploads: Not a must-have for a basic load-testing tool.
- Should-Have: Important but not critical.
    - File Uploads: Should-have for a comprehensive load-testing tool.
- Could-Have: Nice to have but not necessary.
    - File Uploads: Could be seen as a could-have for simpler tools.
- Won't-Have: Not planned for this release.
    - File Uploads: Not applicable.

#### Product Tree
The Product Tree framework involves visualizing features as parts of a tree:
Trunk: Core features of the product.
File Uploads: Not part of the trunk for a basic tool.
Branches: Features available in the next release.
File Uploads: Could be a branch feature for the next release.
Leaves: Future features.
File Uploads: Could be a leaf feature for long-term planning.

### Conclusion
Using these frameworks, we can categorize the file upload feature as follows:

- RICE Score: High priority due to its high reach, impact, and confidence, despite the moderate effort required.
- Kano Model: Performance feature with potential to be an excitement feature for advanced users.
- MoSCoW Method: Should-have feature for a comprehensive tool.
- Product Tree: Branch feature for the next release, with potential to be a leaf feature for future enhancements.


This analysis highlights the importance and prioritization of the file upload feature in an HTTP load-testing and benchmarking library, ensuring it meets user needs and enhances the tool's functionality.