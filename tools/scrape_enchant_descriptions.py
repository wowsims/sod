#!/usr/bin/python

# Usage example:
# python3 ./tools/scrape_enchant_descriptions.py ./tools/database/enchant_overrides.go ./assets/enchants/descriptions.json

import json
import re
import sys

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.chrome.options import Options

if len(sys.argv) < 3:
    raise Exception("Missing arguments, expected input_file_path and output_file_path")
input_file_path = sys.argv[1]
output_file_path = sys.argv[2]

input_file = open(input_file_path, 'r')
input_lines = input_file.readlines()

enchants = []
for line in input_lines:
    spell_id_match = re.search(r"SpellId:\s*(\d+)", line)
    if spell_id_match is None:
        continue
    spell_id = int(spell_id_match.group(1))

    effect_id_match = re.search(r"EffectId:\s*(\d+)", line)
    effect_id = int(effect_id_match.group(1))

    enchants.append({
        "spell_id": spell_id,
        "effect_id": effect_id,
    })

# Added these options so that chrome would run in a docker container
chrome_options = Options()
chrome_options.add_argument("--headless")
chrome_options.add_argument("--no-sandbox")
chrome_options.add_argument("--disable-dev-shm-usage")

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_options)
wait = WebDriverWait(driver, 10)
element_locator = (By.ID, "data-tree-switcher")

def get_spell_effect_description(spell_url):
    driver.get(spell_url)
    wait.until(EC.presence_of_element_located(element_locator))
    details_table = driver.find_elements(By.ID, "spelldetails")[0]
    effect_elem = details_table.find_elements(By.CLASS_NAME, "q2")[0]
    print("Spell {} has description {}".format(spell_url, effect_elem.text))
    return effect_elem.text

def get_enchant_description(enchant):
    return get_spell_effect_description("https://wowhead.com/classic/spell={}".format(enchant["spell_id"]))

for enchant in enchants:
    enchant["description"] = get_enchant_description(enchant)

driver.quit()

with open(output_file_path, "w") as outfile:
    outfile.write("{\n")
    for i, enchant in enumerate(enchants):
        outfile.write("\t\"{}\": \"{}\"{}\n".format(enchant["effect_id"], enchant["description"], "" if i == len(enchants) - 1 else ","))
    outfile.write("}")
