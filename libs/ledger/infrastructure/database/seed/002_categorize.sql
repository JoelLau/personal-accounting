UPDATE entries
SET ledger_accounts_id = 1 -- assets
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    'FAST PAYMENT
OTHR-Other                        to Kaylin Tay Wan%'
    , '%FAST PAYMENT%to Joel IBKR%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;


UPDATE entries
SET ledger_accounts_id = 4101 -- expense:insurance
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    'AXS PTE LTD SINGAPORE SGP'
    , '%INSUR PREM%'
    , '%IBG GIRO%AIA SINGAPORE %'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4601 -- doctor
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    'KKH-KKIVF SINGAPORE SG'
    , 'SPC KKH (OUTPATIENT) - SINGAPORE SG'
  ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;


-- UPDATE entries
-- SET ledger_accounts_id = -- uncategorized income
-- WHERE ledger_accounts_id = 4000
--   AND description LIKE ANY (ARRAY [
--     'BILL PAYMENT - DBS INTERNET/WIRELESS'
--     ])
--   AND debit_microsgd = 0
--   AND credit_microsgd > 0
-- ;

UPDATE entries
SET ledger_accounts_id = 3000 -- uncategorized income
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    'BILL PAYMENT - DBS INTERNET/WIRELESS'
    ])
  AND debit_microsgd = 0
  AND credit_microsgd > 0 -- NOTE: this is CREDIT
;

UPDATE entries
SET ledger_accounts_id = 4201 -- home maintenance
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%MIELE PTE LTD%'
    , '%SELFFIX PTE LTD%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4999 -- misc
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%MANUS AI%'
    , '%KREA.AI%'
    , '%HITEM3D.AI%'
    , '%OPENART AI%'
    , '%ELICIT.COM%'
    ])
;

UPDATE entries
SET ledger_accounts_id = 4305 -- Shopping - K
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%SEPHORA%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4103 -- Parents
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%Frank Tay%'
    , '%CHOW SOUK YONG W%'
    , '%Daisy Chow%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4602 -- Education
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%FIT POUND PTE.%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4102 -- Tax
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%IBG GIRO%ITX%IRAS%TAXS%'
    , 'TANJONG PAGAR TC SINGAPORE SG'
    , '%LOAN PAYMENT%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4204 -- phone
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%WHIZ COMMUNICATIONS %'
    , '%GOMO MOBILE PLAN %'
    , '%PAYPAL *NPARKS%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4203 -- dog
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%VETS FOR PETS %'
    , '%VET PRACTICE %'
    , '%PET LOVERS CENTRE %'
    , 'MYFURBAEBIE SINGAPORE SG'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4304 -- shopping - J
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%AMAZON WEB SERVICES %'
    , '%STEAMGAMES.COM%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4501 -- train (public transport)
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%BUS/MRT%'
    , '%HELLORIDE %'
    , '%WWW.ANYWHEEL.SG%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4301 -- entertainment
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    'NETFLIX.COM SINGAPORE SG'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4401 -- groceries
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%NTUC FP%'
    , '%FAIRPRICE %'
    , '%SHENG SIONG %'
    , 'FINEST FUNAN SINGAPORE SG%'
    , '%DAISO JAPAN%'
    , '%NETS QR%GUARDIAN%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;


UPDATE entries
SET ledger_accounts_id = 4202 -- cleaner
WHERE ledger_accounts_id = 4000
  AND description LIKE '%HELPLING%'
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;



UPDATE entries
SET ledger_accounts_id = 4402 -- eating out
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%THE PALACE KOREAN%'
    , '%40 INDIAN%'
    , '%7-ELEVEN%'
    , '%AAMA BROTHER''S%'
    , '%AGAN GUO KUI%'
    , '%AH LUCK%'
    , '%AJUMMAS %'
    , '%BAHRU PANCAKE %'
    , '%BAHRU PANCAKES %'
    , '%BAN MIAN FISH SOUP%'
    , '%BENGAWAN SOLO %'
    , '%BIANG BIANG NOODLES%'
    , '%BIRDS OF A FEATHER %'
    , '%CAKEBAR %'
    , '%CHAGEE %'
    , '%CHEF HAINANESE WESTERN%'
    , '%CHICHA SAN CHEN %'
    , '%CHICK-FIL-A%'
    , '%Chen Yang%'
    , '%D''LIFE SIGNATURE%'
    , '%DAEBAK KOREAN REST%'
    , '%DON DON DONKI %'
    , '%FAROK NISA FAMILY FOOD%'
    , '%FAST PAYMENT% to Qashier-%'
    , '%FAST PAYMENT%to Yati%'
    , '%FAST PAYMENT%to Brian Rrrrrr%'
    , '%FRUIT BOX %'
    , '%FUND TRANSFER%FOMO PAY%'
    , '%GENKI SUSHI%'
    , '%GLAM (85) PTE LTD %'
    , '%GNS FOODS PTE.%'
    , '%GOCHISO SHOKUDO PTE%'
    , '%GOLDMOON SUNTEC%'
    , '%GYG -%'
    , '%HANBAOBAO PTE LTD%'
    , '%HELLOFISH STEAM%'
    , '%HEY CANTON PTE.%'
    , '%HEYTEA%'
    , '%HOMIES BAKER%'
    , '%HTI3 PTE LTD%'
    , '%JIU MAO SARAWAK KOLO MEE%'
    , '%JOY STEAKS%'
    , '%JOYFUL CCK%'
    , '%KEBABS FAKTORY %'
    , '%KOO KEE %'
    , '%KOPITIAM %'
    , '%KOPITIAM %'
    , '%KOUFU %'
    , '%KUMO JAPANESE DINING %'
    , '%LADERACH SWISS CHOC%'
    , '%LEROY LIM YING R%'
    , '%LIHO %'
    , '%LILI KWAY CHAP BRAISED DUCK DUCK%'
    , '%LITTLE FISHER %'
    , '%LONG JOHN SILVER%'
    , '%MAGURO BROTHERS %'
    , '%MANDY CHUA LI MI%'
    , '%MASTER PRATA %'
    , '%MCDONALD%'
    , '%MR AVOCADO%'
    , '%MR.K WESTERN KIT%'
    , '%MUCHACHOS%'
    , '%NAM KEE %'
    , '%NETS QR%85 DASNP%'
    , '%NETS QR%RIVERSIDE%'
    , '%OLD CHANG KEE%'
    , '%OMMAS%'
    , '%OMNIVORE %'
    , '%OTTIE PANCAKES %'
    , '%PASTAGO %'
    , '%PETREL CATERING%'
    , '%PaperBakes %'
    , '%QB NET INTERNATI%'
    , '%RAZER MERCHANT%'
    , '%RUMEL %'
    , '%SAGYE KOREAN '
    , '%SAGYE KOREAN RESTAURAN%'
    , '%SALAD STOP %'
    , '%SF - ONE RAFFLES SINGAPORE SG%'
    , '%SF - VIVOCITY SINGAPORE SG%'
    , '%SHIHLIN TAIWAN%'
    , '%SHOCK BURGER%'
    , '%SIN KEE CHICKEN RICE%'
    , '%SMP*LI XING FOODSTUFF%'
    , '%SO GOOD BAKERY - AMARA SINGAPORE SG%'
    , '%SP PAPERPALETTE SINGAPORE SG%'
    , '%STUFF''D %'
    , '%SUBWAY @%'
    , '%SUKIYA %'
    , '%SUPER SUSHI %'
    , '%SUSHI EXPRESS%'
    , '%SUSHI ZUSHI%'
    , '%TAKAGI RAMEN%'
    , '%THAI DESSERT%'
    , '%THE ALLEY %'
    , '%THE BACKYARD BAKERS RP SINGAPORE SG%'
    , '%THE DAILY CUT %'
    , '%THE GOOD PASTA %'
    , '%THE SOUP SPOON PTE LTD%'
    , '%TIAN XIANG WANTON MEE%'
    , '%TIONG BAHRU PAU%'
    , '%TODAY RESTAURANT SINGAPORE SG%'
    , '%TONGUE TIP LZ BEEF NOO SINGAPORE SG%'
    , '%TORA JAPANESE BBQ_F SINGAPORE SG%'
    , '%TORI-Q%'
    , '%TWO MEN BAGEL HOUSE %'
    , '%WHOLE EARTH %'
    , '%WOK HEY %'
    , '%XIANG XIANG %'
    , '%XIN YUAN JI %'
    , '%XW PLUS WESTERN GRILL %'
    , '%XW PLUS WESTERN GRILL%'
    , '%YA KUN KAYA TOAST%'
    , 'DA XI-GUOCCO TOWER SINGAPORE SG'
    , 'SEOUL BUNSIK SINGAPORE SG'
    , '%CHENMAPO %'
    , '%BREADTALK %'
    , '%VISTA PANCAKES %'
    , '%JOY LUCK TEAHOUSE%'
    , '%HOCKHUA %'
    , '%LIU LANG MIAN %'
    , '%SWEE HUAT %'
])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;
