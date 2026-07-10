#!/usr/bin/env ruby
# frozen_string_literal: true

require "cgi"
require "json"

root = File.expand_path("..", __dir__)
reader_path = File.join(root, "docs/capstone/index.html")
reader = File.read(reader_path)
errors = []
doc_data = reader[%r{<script id="doc-data"[^>]*>(.*?)</script>}m, 1]
documents = doc_data ? JSON.parse(doc_data) : []
errors << "expected 9 doc-data entries, found #{documents.length}" unless documents.length == 9

documents.each do |document|
  relative_path = document.fetch("rel")
  marker_offset = reader.index(%(id="#{document.fetch("id")}"))
  unless marker_offset
    errors << "missing embedded document: #{relative_path}"
    next
  end

  article_start = reader.rindex("<article", marker_offset)
  article_end = reader.index("</article>", marker_offset)
  article = reader[article_start...article_end]
  unless article.include?(%(data-path="#{relative_path}"))
    errors << "embedded document path differs: #{relative_path}"
    next
  end
  embedded = article&.match(%r{<div class="source">\s*<pre>(.*?)</pre}m)&.[](1)
  unless embedded
    errors << "missing embedded source: #{relative_path}"
    next
  end

  embedded = CGI.unescapeHTML(embedded).sub(/\A\n/, "")
  source = File.read(File.join(root, relative_path))
  errors << "embedded source differs: #{relative_path}" unless embedded == source

  title = document.fetch("title")
  search = article.match(/data-search=(.)(.*?)\1/m)&.[](2)
  expected_search = "#{title.downcase} #{relative_path.downcase} #{source.downcase}"
  unless search && CGI.unescapeHTML(search) == expected_search
    errors << "embedded search differs: #{relative_path}"
  end
end

ids = reader.scan(/\bid="([^"]+)"/).flatten.to_h { |id| [id, true] }
reader.scan(/href="#([^"]+)"/).flatten.each do |fragment|
  errors << "broken fragment: ##{fragment}" unless ids[fragment]
end

reader.scan(/href="([^"]+)"/).flatten.uniq.each do |href|
  next if href.empty? || href.start_with?("#", "http://", "https://", "mailto:")

  target = href.split("#", 2).first
  path = File.expand_path(target, File.dirname(reader_path))
  errors << "missing local target: #{href}" unless File.exist?(path)
end

if errors.empty?
  puts "capstone-reader: embedded sources and links are valid"
  exit 0
end

errors.each { |error| warn "capstone-reader: #{error}" }
exit 1
