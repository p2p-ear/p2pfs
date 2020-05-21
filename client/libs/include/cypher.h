class XORcypher {
private:
    char secret_[4096];
public:
    // Initialize secret with random bytes
    XORcypher();

    // Initialize secret with given bytes
    XORcypher(char* secret);

    // Retrieve secret
    const char* get_secret();

    // Cypher shard in-place
    void operator()(char* shard);
};