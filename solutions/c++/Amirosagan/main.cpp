#include <boost/multiprecision/cpp_int.hpp>
#include <chrono>
#include <fstream>
#include <iostream>
#include <string>

struct City {
  std::string name;
  boost::multiprecision::int128_t min = 1000000000;
  boost::multiprecision::int128_t max = 0;
  boost::multiprecision::int128_t sum = 0;
  boost::multiprecision::int128_t count = 0;
  double avg = 0;
  boost::multiprecision::int128_t temperature = 0;
};

const int MOD = 100007;

int main() {
  // trace time
  std::chrono::steady_clock::time_point begin =
      std::chrono::steady_clock::now();

  std::cout << "working on it..." << std::endl;
  std::ifstream file("../../../generation-binaries/data.txt");

  if (!file) {
    std::cout << "Error opening file" << std::endl;
    return 1;
  }
  char ch;
  std::string city;
  std::string population;
  std::string line;
  bool iscity = true;
  long long cityHash = 0;

  City *umap = new City[MOD];

  while (file.get(ch)) {
    if (ch == ';') {
      population = line;
      boost::multiprecision::int128_t *temp =
          new boost::multiprecision::int128_t(
              boost::multiprecision::cpp_int(population));
      umap[cityHash].temperature = *temp;
      umap[cityHash].count += 1;
      umap[cityHash].sum += *temp;
      umap[cityHash].min =
          umap[cityHash].min > *temp ? *temp : umap[cityHash].min;
      umap[cityHash].max =
          umap[cityHash].max < *temp ? *temp : umap[cityHash].max;
      umap[cityHash].avg =
          (double)umap[cityHash].sum / (double)umap[cityHash].count;
      umap[cityHash].name = city;

      line.clear();
      city.clear();
      population.clear();
      cityHash = 0;
      iscity = true;
    } else if (ch == ',') {
      city = line;
      line.clear();
      iscity = false;
    } else {
      if (iscity) {
        cityHash = ((cityHash * 31) % MOD + ch) % MOD;
      }
      line += ch;
    }
  }

  for (int i = 0; i < MOD; i++) {
    if (umap[i].count != 0) {
      std::cout << umap[i].name << " " << umap[i].min << " " << umap[i].max
                << " " << umap[i].avg << std::endl;
    }
  }

  file.close();

  // trace time

  std::ofstream out("trace_output.trace");

  for (int i = 0; i < MOD; i++) {
    if (umap[i].count != 0) {
      out << umap[i].name << " " << umap[i].min << " " << umap[i].max << " "
          << umap[i].avg << std::endl;
    }
  }

  std::chrono::steady_clock::time_point end = std::chrono::steady_clock::now();

  out << "Time = "
      << std::chrono::duration_cast<std::chrono::seconds>(end - begin).count()
      << "[s]" << std::endl;

  out << "Time = "
      << std::chrono::duration_cast<std::chrono::microseconds>(end - begin)
             .count()
      << "[µs]" << std::endl;

  return 0;
}
