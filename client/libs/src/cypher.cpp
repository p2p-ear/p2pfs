#include "../include/cypher.h"
#include <random>
#include <cstring>

XORcypher::XORcypher() {
    std::independent_bits_engine<std::default_random_engine, __CHAR_BIT__, char> engine;
    for (unsigned int i = 0; i < 4096; i++) {
        secret_[i] = engine();
    }
}

XORcypher::XORcypher(char* secret) {
    std::memcpy(secret_, secret, 4096);
}

const char* XORcypher::get_secret() {
    return secret_;
}

void XORcypher::operator()(char* shard) {
    for (unsigned int i = 0; i < 4096; i++) {
        shard[i] ^= secret_[i];
    }
}