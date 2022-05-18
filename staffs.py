#!/usr/bin/env python3

import sys, os
import argparse
import urllib.request
import urllib.parse
from bs4 import BeautifulSoup

cName = os.path.basename(sys.argv[0])
cVersion = '1.0.0'
cCopyright = "Copyright (C) by AD, IX 2021. All rights reserved."


urlEndpoint = r"https://staffsearch.fernuni-hagen.de/"


def printDataset (*, data, headers=None):
  if headers is None:
    headers = ['•' for d in data]
    sep = ' '
  else:
    sep = ': '
  for i in range(0, len(data)):
    d = str(headers[i]).rjust(10) if headers[i] is not None else ''.rjust(10)
    print (f"{d}{sep}{data[i]!s}")


def printDatatab (*, dataArr, headers=None):
  if headers is not None:
    for d in headers:
      print (f"{d:12} ", end='')
    print("")
    for data in dataArr:
      for d in data:
        print (f"{d:12} ", end='')
      print("")


def preprocessData (data):
  try:
    data = data[data.index(b'<p id="result"'):]
  except ValueError:
    data = b''
  data = data[0:data.find(b'<!--')]
  return b"<root>\n" + data + b"</root>\n"


def queryPerson (searchterm, *, show_raw=False):
  query = { 's': searchterm, 'p': "1" }
  strQuery = f"?{urllib.parse.urlencode(query)}"

  hdrUA = r"Mozilla/5.0 (X11; U; Linux i686) Gecko/20071127 Firefox/2.0.0.11"
  headers = { 'User-Agent': hdrUA }
  request = urllib.request.Request(urlEndpoint + strQuery, data=None, headers=headers)

  with urllib.request.urlopen(request) as response:
    #print (f"Response Status: {response.status}\n\n")

    data = preprocessData(response.read())
    soup = BeautifulSoup(data, 'html.parser')

    eltTable = soup.find('table')
    if eltTable is None:
      print(f"- No result returned for ”{searchterm}“")
      return -1

    eltsTR = eltTable.find_all('tr')

    dataHeaders = [e.get_text() for e in eltsTR[0].find_all('th')]
    if show_raw: print(dataHeaders)

    dataResults = []
    for tr in eltsTR[1:]:
      dataResults.append([e.get_text() for e in tr.find_all('td')])
    if show_raw:
      for result in dataResults:
        print(result)

    for data in dataResults:
      printDataset(data=data, headers=dataHeaders)
      print('')

    #printDatatab(dataArr=dataResults, headers=dataHeaders)

    return 0


def BannerFormatterFactory (bannertext=None, copyright=None, formatter_class=argparse.HelpFormatter):

  class _Formatter(formatter_class):
    def _format_header (self):
      nonlocal bannertext, copyright
      if bannertext is None: bannertext = f"{self._prog}"
      if copyright is None: copyright = f"Copyright (C) by the respective Author. All rights reserved."
      banner = ''
      if bannertext != '': banner += f"{bannertext}\n"
      if copyright != '': banner += f"{copyright}\n"
      if bannertext is not None or copyright is not None: banner += "\n"
      return banner

    def add_usage (self, *args):
      self._add_item(self._format_header, [])
      super().add_usage(*args)

  return _Formatter


if __name__ == "__main__":

  #parser = argparse.ArgumentParser(description='Query the StaffSearch service of the FernUni.', formatter_class=MyFormatter)
  parser = argparse.ArgumentParser(
            description='Query the StaffSearch service of the FernUniversität in Hagen.',
            formatter_class=BannerFormatterFactory(bannertext=cName+" – Perform query using StaffSearch",
                                                   copyright=cCopyright))
  parser.add_argument('search_terms',
                      type=str, nargs='+',
                      help='the query terms')
  parser.add_argument('-e', '--endpoint', metavar='endpoint',
                      type=str, nargs=1, default=urlEndpoint,
                      help='StaffSearch endpoint (default: %(default)s)')
  parser.add_argument('-r', '--raw',
                      action='store_true',
                      help='show raw results')
  parser.add_argument('-V', '--version', action='version',
                      version='%(prog)s '+cVersion,
                      help='the version')
  args = parser.parse_args()

  print('')
  for query in args.search_terms:
    queryPerson(query, show_raw=args.raw)
    print('')
