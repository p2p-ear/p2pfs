#include "../include/cypher.h"
#include <random>
#include <cstring>

XORcypher::XORcypher(unsigned int size) : size_(size), secret_(new char[size]) {
    std::independent_bits_engine<std::default_random_engine, __CHAR_BIT__, char> engine;
    for (unsigned int i = 0; i < size_; i++) {
        secret_[i] = engine();
    }
}

XORcypher::XORcypher(unsigned int size, char* secret) : size_(size), secret_(new char[size]) {
    std::memcpy(secret_, secret, size_);
}

XORcypher::~XORcypher() {
    delete[] secret_;
}

const char* XORcypher::get_secret() {
    return secret_;
}

const int XORcypher::get_size() {
    return size_;
}

void XORcypher::operator()(char* shard) {
    for (unsigned int i = 0; i < size_; i++) {
        shard[i] ^= secret_[i];
    }
}

void XORcypher::operator()(char* shard, const unsigned int shardsize) {
    for (unsigned int i = 0; i < shardsize; i++) {
        shard[i] ^= secret_[i % size_];
    }
}