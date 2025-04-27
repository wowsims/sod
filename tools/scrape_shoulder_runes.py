#!/usr/bin/python

# This tool generates the classic SoD runes data

import sys
import requests
import math
import re

from typing import List

from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.chrome.options import Options

if len(sys.argv) < 2:
    raise Exception("Missing arguments, expected output_file_path")

WowheadBranch = "classic"
WowheadBranchPTR = "classic-ptr"

CURRENT_BRANCH = WowheadBranchPTR

output_file_path = sys.argv[1]

# Added these options so that chrome would run in a docker container
chrome_options = Options()
chrome_options.add_argument("--headless")
chrome_options.add_argument("--no-sandbox")
chrome_options.add_argument("--disable-dev-shm-usage")

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_options)
wait = WebDriverWait(driver, 10)
element_locator = (By.ID, "data-tree-switcher")

def _get_id_from_link(link):
    return int(link.split("/")[-2].split("=")[-1])


def get_item_ids() -> List[int]:
    driver.get(f"https://www.wowhead.com/{CURRENT_BRANCH}/items/name:%22Soul+of+the%22?filter=82:142;2:0;11506:spell_holy_divinespirit")
    wait.until(EC.presence_of_element_located(element_locator))

    listview = driver.find_element(By.ID, "lv-items")
    pages = int(listview.find_element(By.CLASS_NAME, "listview-nav").find_element(By.CSS_SELECTOR, 'b:last-child').text)/50
    pages = math.ceil(pages)
    all_ids = []

    for page in range(pages):
        print(f'Loading page {page} for runes...')
        driver.get(f"https://www.wowhead.com/{CURRENT_BRANCH}/items/name:%22Soul+of+the%22?filter=82:142;2:0;11506:spell_holy_divinespirit#{page*50}")
        driver.refresh()
        wait.until(EC.presence_of_element_located(element_locator))
        listview = driver.find_element(By.ID, "lv-items")
        rows = listview.find_elements(By.CLASS_NAME, "listview-row")
        all_ids.extend([_get_id_from_link(row.find_element(By.CLASS_NAME, "listview-cleartext").get_attribute("href"))
            for row in rows])

    return all_ids

def get_tooltips_response(id):
    # Get the underlying item ID from the engraving ID
    url = f"https://nether.wowhead.com/{CURRENT_BRANCH}/tooltip/item/{id}"
    result = requests.get(url)

    if result.status_code == 200:
        response_json = result.text
        return response_json
    else:
        print(f"Request for id {id} failed with status code: {result.status_code}")
        return None
    

# id, tooltip_json
to_export = []

item_ids = get_item_ids()

print(f"Export Count ({len(item_ids)}) {item_ids}")

to_export = []

mismatchedIds = {
    "1219819": 1220096, # Soul of Echoes
    "1219947": 1220355, # Soul of the Night
    "1219740": 1219957, # Soul of the Tactician
    # Soul of the Refined has 4 versions - Paladin, Priest, Warlock, Shaman and the scraper can't correctly separate them
    "1219870": 1220198, # Paladin
    "1219825": 1220108, # Priest
    "1219903": 1220264, # Shaman
    "1219786": 1220034, # Warlock
    "1233551": 1220352  # Soul of the Starcaller
}

for item_id in item_ids:
    item_response = get_tooltips_response(item_id) or ""
    if not item_response:
        continue
    spell_ids = re.findall(r'\/spell=(\d+)', item_response)

    if len(spell_ids) >= 2:
        # The base spell is different from what we typically use but the spell we actually want appears as the first related spell in the "See also" tab
        enchant_spell_id = spell_ids[0]

        if enchant_spell_id in mismatchedIds:
            to_export.append([mismatchedIds[enchant_spell_id], item_response]) 
        else:
            try:
                driver.get(f"https://www.wowhead.com/{CURRENT_BRANCH}/spell={enchant_spell_id}#see-also-other")
                driver.refresh()
                wait.until(EC.presence_of_element_located(element_locator))

                see_also_tab = driver.find_element(By.ID, "tab-see-also-other")
                rows = see_also_tab.find_elements(By.CLASS_NAME, "listview-row")

                if len(rows) > 0:
                    actual_spell_id = _get_id_from_link(rows[0].find_element(By.CLASS_NAME, "listview-cleartext").get_attribute("href"))
                    # Use the Spell ID of the effect but keep the item tooltip for class matching and display
                    to_export.append([actual_spell_id, item_response])
                else:
                    print(f"No related spell IDs found for spell ID {enchant_spell_id}")
            except NoSuchElementException:
                print(f"No related spell IDs found for spell ID {enchant_spell_id}")
    else:
        print(f"Less than 2 spell IDs found for item ID {id}")

driver.quit()
output_string = '\n'.join([str(','.join([str(inner_elem) for inner_elem in elem])) for elem in to_export])

with open(output_file_path, "w") as outfile:
    outfile.write(output_string)
