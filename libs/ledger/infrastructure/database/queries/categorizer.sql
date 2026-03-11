UPDATE entries
SET ledger_accounts_id = 1 -- assets
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    'FAST PAYMENT
OTHR-Other                        to Kaylin Tay Wan%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

--

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
  AND debit_microsgd > 0
  AND credit_microsgd = 0
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
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;

UPDATE entries
SET ledger_accounts_id = 4304 -- shopping - j
WHERE ledger_accounts_id = 4000
  AND description LIKE ANY (ARRAY [
    '%AMAZON WEB SERVICES %'
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
    , '%GYG -%'
    , '%SHIHLIN TAIWAN%'
    , '%TORI-Q%'
    , '%CHICK-FIL-A%'
    , '%GOCHISO SHOKUDO PTE%'
    , '%SUSHI ZUSHI%'
    , '%OLD CHANG KEE%'
    , '%MCDONALD%'
    , '%YA KUN KAYA TOAST%'
    , '%SMP*LI XING FOODSTUFF%'
    , '%HEYTEA%'
    , '%LONG JOHN SILVER%'
    , '%MUCHACHOS%'
    , '%STUFF''D %'
    , '%BAHRU PANCAKE %'
    , '%KEBABS FAKTORY %'
    , '%HOMIES BAKER%'
    , '%TWO MEN BAGEL HOUSE %'
    , '%PASTAGO %'
    , '%CAKEBAR %'
    , '%KUMO JAPANESE DINING %'
    , '%MASTER PRATA %'
    , '%BIRDS OF A FEATHER %'
    , '%XW PLUS WESTERN GRILL %'
    , '%KOPITIAM %'
    , '%XIN YUAN JI %'
    , '%AJUMMAS %'
    , '%SUKIYA %'
    , '%OMNIVORE %'
    , '%SUPER SUSHI %'
    , '%FRUIT BOX %'
    , '%CHICHA SAN CHEN %'
    , '%THE DAILY CUT %'
    , '%KOPITIAM %'
    , '%KOUFU %'
    , '%WHOLE EARTH %'
    , '%XIANG XIANG %'
    , '%SAGYE KOREAN '
    , '%NAM KEE %'
    , '%WOK HEY %'
    , '%SALAD STOP %'
    , '%KOO KEE %'
    , '%RUMEL %'
    , '%7-ELEVEN%'
    , '%BENGAWAN SOLO %'
    , '%BIANG BIANG NOODLES%'
    , '%CHAGEE %'
    , '%THE ALLEY %'
    , '%BAHRU PANCAKES %'
    , '%PaperBakes %'
    , '%DON DON DONKI %'
    , '%THE GOOD PASTA %'
    , '%MR.K WESTERN KIT%'
    , '%XW PLUS WESTERN GRILL%'
    , '%TORA JAPANESE BBQ_F SINGAPORE SG%'
    , '%TONGUE TIP LZ BEEF NOO SINGAPORE SG%'
    , '%TODAY RESTAURANT SINGAPORE SG%'
    , '%THE BACKYARD BAKERS RP SINGAPORE SG%'
    , '%SP PAPERPALETTE SINGAPORE SG%'
    , '%SO GOOD BAKERY - AMARA SINGAPORE SG%'
    , '%SHOCK BURGER%'
    , '%DAEBAK KOREAN RESTG%'
    , '%SF - VIVOCITY SINGAPORE SG%'
    , '%SF - ONE RAFFLES SINGAPORE SG%'
    , '%SAGYE KOREAN RESTAURAN%'
    , '%TIONG BAHRU PAU%'
    , '%SIN KEE CHICKEN RICE%'
    , '%NETS QR%RIVERSIDE%'
    , '%MR AVOCADO%'
    , '%LILI KWAY CHAP BRAISED DUCK DUCK%'
    , '%JOY STEAKS%'
    , '%JIU MAO SARAWAK KOLO MEE%'
    , '%FAST PAYMENT% to Qashier-%'
    , '%HANBAOBAO PTE LTD%'
    , '%FAROK NISA FAMILY FOOD%'
    , '%GLAM (85) PTE LTD %'
    , '%AGAN GUO KUI%'
    , '%AH LUCK%'
    , '%AAMA BROTHER''S%'
    , '%NETS QR%85 DASNP%'
    , '%HELLOFISH STEAM%'
    , '%TAKAGI RAMEN%'
    , '%SUSHI EXPRESS%'
    , '%BAN MIAN FISH SOUP%'
    , '%MAGURO BROTHERS %'
    , '%LITTLE FISHER %'
    , '%LADERACH SWISS CHOC%'
    , '%Chen Yang%'
    , '%RAZER MERCHANT%'
    , '%QB NET INTERNATI%'
    , '%LIHO %'
    , '%MANDY CHUA LI MI%'
    , '%LEROY LIM YING R%'
    , '%GNS FOODS PTE.%'
    , '%HEY CANTON PTE.%'
    , '%GOLDMOON SUNTEC%'
    , '%JOYFUL CCK%'
    , '%FUND TRANSFER%FOMO PAY%'
    , '%D''LIFE SIGNATURE%'
    ])
  AND debit_microsgd > 0
  AND credit_microsgd = 0
;


