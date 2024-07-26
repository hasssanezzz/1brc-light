#include <boost/multiprecision/cpp_int.hpp>
#include <chrono>
#include <fstream>
#include <iostream>
#include <string>
#include <unordered_map>

struct City {
  boost::multiprecision::int128_t min = 1000000000;
  boost::multiprecision::int128_t max = 0;
  boost::multiprecision::int128_t sum = 0;
  boost::multiprecision::int128_t count = 0;
  double avg = 0;
  boost::multiprecision::int128_t temperature = 0;
};

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

  std::unordered_map<std::string, City> umap;

  while (file.get(ch)) {
    if (ch == ';') {
      population = line;
      boost::multiprecision::int128_t *temp =
          new boost::multiprecision::int128_t(
              boost::multiprecision::cpp_int(population));
      umap[city].temperature = *temp;
      umap[city].count += 1;
      umap[city].sum += *temp;
      umap[city].min = umap[city].min > *temp ? *temp : umap[city].min;
      umap[city].max = umap[city].max < *temp ? *temp : umap[city].max;
      umap[city].avg = (double)umap[city].sum / (double)umap[city].count;

      line.clear();
      city.clear();
      population.clear();
    } else if (ch == ',')
      city = line, line.clear();
    else
      line += ch;
  }

  for (const auto &pair : umap) {
    std::cout << pair.first << " " << pair.second.min << " " << pair.second.max
              << " " << pair.second.avg << std::endl;
  }

  file.close();

  // trace time

  std::ofstream out("trace_output.trace");

  for (const auto &pair : umap) {
    out << pair.first << " " << pair.second.min << " " << pair.second.max << " "
        << pair.second.avg << std::endl;
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
