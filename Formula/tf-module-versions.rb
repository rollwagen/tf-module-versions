# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class TfModuleVersions < Formula
  desc ""
  homepage ""
  version "0.0.4"
  license "Apache-2.0"

  on_macos do
    url "https://github.com/rollwagen/tf-module-versions/releases/download/v0.0.4/tf-module-versions_0.0.4_darwin_all.tar.gz"
    sha256 "0943770434eb67001c92f0b7004844f487dcc89f8457a5df6d7d3dde66ecad7b"

    def install
      bin.install "tf-module-versions"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/rollwagen/tf-module-versions/releases/download/v0.0.4/tf-module-versions_0.0.4_linux_arm64.tar.gz"
      sha256 "a4e9dab78bf86f25ef5ae0e1d517531cb2584ee3c485eaf7cf90610dd7c16e2e"

      def install
        bin.install "tf-module-versions"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/rollwagen/tf-module-versions/releases/download/v0.0.4/tf-module-versions_0.0.4_linux_amd64.tar.gz"
      sha256 "5bec45d0673eba12dcd88da645be947135c06893b7dd16c0c2c7c5983ebef7bf"

      def install
        bin.install "tf-module-versions"
      end
    end
  end
end
