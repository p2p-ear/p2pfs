class XORcypher {
private:
    char* secret_;
    const unsigned int size_;
public:
    // Initialize secret with random bytes
    XORcypher(unsigned int size);

    // Initialize secret with given bytes
    XORcypher(unsigned int size, char* secret);

    ~XORcypher();

    // Retrieve secret
    const char* get_secret();

    // Retrieve secret size
    const int get_size();

    // Cypher shard in-place
    void operator()(char* shard);
};